package linkaudit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/archive"
)

type ContentOptions struct {
	ContentDir, Report string
	Posts              []string
	Resume, UseModel   bool
	Delay              time.Duration
	Model              archive.Options
}

type ContentResult struct {
	Post     string   `json:"post"`
	HasIssue bool     `json:"has_issue"`
	Issues   []string `json:"issues"`
	Error    string   `json:"error,omitempty"`
}

var wpShortcode = regexp.MustCompile(`(?i)\[/?(?:gallery|embed|caption|video|audio|playlist)(?:\s[^\]\n]*)?\]`)
var htmlCode = regexp.MustCompile(`(?is)<(pre|code)\b[^>]*>.*?</(?:pre|code)\s*>`)
var hugoCode = regexp.MustCompile(`(?s)\{\{[<%].*?[>%]\}\}`)
var inlineCode = regexp.MustCompile("`+[^`\n]*`+")

func mask(source string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' {
			return r
		}
		return ' '
	}, source)
}

// WordPressIssues ignores examples in front matter, code and Hugo shortcodes,
// and distinguishes shortcodes from ordinary Markdown links and references.
func WordPressIssues(content string) []string {
	lines := strings.Split(content, "\n")
	front := ""
	fence := ""
	fenceLen := 0
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if i == 0 && (trim == "---" || trim == "+++") {
			front = trim
			lines[i] = mask(line)
			continue
		}
		if front != "" {
			if trim == front {
				front = ""
			}
			lines[i] = mask(line)
			continue
		}
		if fence != "" {
			if strings.HasPrefix(trim, strings.Repeat(fence, fenceLen)) && strings.Trim(trim, fence+" \t") == "" {
				fence = ""
			}
			lines[i] = mask(line)
			continue
		}
		if strings.HasPrefix(trim, "```") || strings.HasPrefix(trim, "~~~") {
			fence = trim[:1]
			fenceLen = len(trim) - len(strings.TrimLeft(trim, fence))
			lines[i] = mask(line)
			continue
		}
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			lines[i] = mask(line)
		}
	}
	content = strings.Join(lines, "\n")
	content = htmlCode.ReplaceAllStringFunc(content, mask)
	content = hugoCode.ReplaceAllStringFunc(content, mask)
	content = inlineCode.ReplaceAllStringFunc(content, mask)
	issues := []string{}
	for _, loc := range wpShortcode.FindAllStringIndex(content, -1) {
		start, end := loc[0], loc[1]
		escapes := 0
		for j := start - 1; j >= 0 && content[j] == '\\'; j-- {
			escapes++
		}
		if escapes%2 == 1 {
			continue
		}
		after := strings.TrimLeft(content[end:], " \t")
		if strings.HasPrefix(after, "(") || strings.HasPrefix(after, "[") || strings.HasPrefix(after, ":") {
			continue
		}
		issues = append(issues, fmt.Sprintf("WordPress shortcode %s at line %d", content[start:end], strings.Count(content[:start], "\n")+1))
	}
	return issues
}

const contentPrompt = `/no_think
Audit this Hugo Markdown post for WordPress migration issues ONLY.
Flag only WordPress square-bracket shortcodes: [gallery], [embed], [caption], [video], [audio], [playlist], including their attributes and closing tags.
Do not flag Markdown links or references, HTML tags, Hugo shortcodes, front matter, Markdown formatting, or code blocks/inline code. Be extremely conservative.
Return ONLY a JSON object: {"has_issue":false,"issues":[]} with issues as an array of strings.
Post: `

func AuditContent(ctx context.Context, o ContentOptions, log io.Writer) ([]ContentResult, error) {
	if log == nil {
		log = io.Discard
	}
	if o.Report == "" || o.Delay < 0 {
		return nil, fmt.Errorf("invalid content audit options")
	}
	if err := o.Model.Validate(); err != nil {
		return nil, err
	}
	names, err := bundles(o.ContentDir, o.Posts)
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("no post bundles found")
	}
	unlock, err := archive.Lock(ctx, o.Report+".lock", o.Model.LockTimeout)
	if err != nil {
		return nil, err
	}
	defer unlock()
	results := []ContentResult{}
	if data, err := os.ReadFile(o.Report); err == nil {
		if err := json.Unmarshal(data, &results); err != nil {
			return nil, err
		}
		if results == nil {
			return nil, fmt.Errorf("content audit report must be an array")
		}
		for _, r := range results {
			if r.Post == "" {
				return nil, fmt.Errorf("content audit record has no post")
			}
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	indices := map[string]int{}
	for i, r := range results {
		indices[r.Post] = i
	}
	v := archive.New(o.Model, log)
	processed := false
	for _, name := range names {
		if err := ctx.Err(); err != nil {
			return results, err
		}
		if _, ok := indices[name]; ok && o.Resume {
			continue
		}
		if processed {
			if err := wait(ctx, o.Delay); err != nil {
				return results, err
			}
		}
		processed = true
		file, err := archive.PostPath(o.ContentDir, "/p/"+url.PathEscape(name)+"/")
		if err != nil {
			return results, err
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return results, err
		}
		r := ContentResult{Post: name, Issues: []string{}}
		if o.UseModel {
			text := []rune(string(data))
			if len(text) > o.Model.MaxContent {
				text = text[:o.Model.MaxContent]
			}
			analysis, err := v.AnalyzeContentJSON(ctx, contentPrompt+name+"\n\n"+string(text))
			if ctx.Err() != nil {
				return results, ctx.Err()
			}
			if err == nil {
				r.HasIssue, _ = analysis["has_issue"].(bool)
				issues, ok := analysis["issues"].([]any)
				if !ok {
					err = fmt.Errorf("model response missing issues array")
				} else {
					for _, issue := range issues {
						s, ok := issue.(string)
						if !ok {
							err = fmt.Errorf("model issue must be a string")
							break
						}
						r.Issues = append(r.Issues, s)
					}
				}
			}
			if err != nil {
				r.Error = err.Error()
				r.HasIssue = true
				r.Issues = []string{"Error analyzing post: " + err.Error()}
			}
		} else {
			r.Issues = WordPressIssues(string(data))
			r.HasIssue = len(r.Issues) > 0
		}
		if i, ok := indices[name]; ok {
			results[i] = r
		} else {
			indices[name] = len(results)
			results = append(results, r)
		}
		fmt.Fprintf(log, "%s: %d potential issues\n", name, len(r.Issues))
		if err := archive.SaveJSON(o.Report, results); err != nil {
			return results, err
		}
	}
	return results, nil
}
