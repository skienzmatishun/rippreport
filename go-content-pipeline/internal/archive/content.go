package archive

import (
	"net/url"
	"path"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func siteHost(host string) bool  { return host == "rippreport.com" || host == "www.rippreport.com" }
func hasExtension(p string) bool { return path.Ext(strings.TrimRight(p, "/")) != "" }
func isPDF(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && strings.EqualFold(path.Ext(u.Path), ".pdf")
}
func tagName(p string) string {
	parts := strings.Split(strings.Trim(p, "/"), "/")
	if len(parts) == 2 && parts[0] == "tags" {
		return parts[1]
	}
	return ""
}
func missingSlug(p string) string {
	s := strings.Trim(p, "/")
	if s == "" || strings.Contains(s, "/") {
		return ""
	}
	switch s {
	case "p", "tags", "categories", "about", "contact", "search", "index.xml", "index.html":
		return ""
	}
	return s
}

type document struct {
	Text                           string
	Paragraphs, Headings, Articles int
	Links                          []string
}

// Parse HTML as a tree so nested Wayback banners cannot count as article text.
func inspectHTML(source string) document {
	root, err := html.Parse(strings.NewReader(source))
	if err != nil {
		return document{}
	}
	var d document
	var text strings.Builder
	var visit func(*html.Node)
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "meta", "link", "noscript", "template":
				return
			}
			for _, a := range n.Attr {
				v := strings.ToLower(a.Val)
				if (a.Key == "id" || a.Key == "class") && (strings.Contains(v, "wayback") || v == "wm-ipp-base" || v == "wm-ipp") {
					return
				}
				if n.Data == "a" && a.Key == "href" {
					d.Links = append(d.Links, a.Val)
				}
			}
			switch n.Data {
			case "p":
				d.Paragraphs++
			case "h1", "h2", "h3", "h4", "h5", "h6":
				d.Headings++
			case "article", "main":
				d.Articles++
			}
		}
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
			text.WriteByte(' ')
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(root)
	d.Text = strings.Join(strings.Fields(text.String()), " ")
	return d
}

func (d document) substantial(paywall bool) bool {
	t := strings.ToLower(d.Text)
	for _, phrase := range []string{"page not found", "access denied", "403 forbidden", "this page is not archived", "no archived versions", "wayback machine doesn't have that page", "checking your browser", "verify you are human"} {
		if strings.Contains(t, phrase) {
			return false
		}
	}
	if !paywall {
		return len([]rune(d.Text)) >= 200 && (d.Paragraphs > 0 || d.Headings > 0)
	}
	count := 0
	for _, phrase := range []string{"subscribe", "subscription", "sign in", "log in", "create account", "register", "paywall", "premium content", "members only"} {
		if strings.Contains(t, phrase) {
			count++
		}
	}
	return count < 3 && len([]rune(d.Text)) >= 500 && (d.Articles > 0 || (d.Paragraphs > 3 && d.Headings > 0))
}

// ReplaceURLs substitutes exact destinations without reprocessing newly inserted
// archive URLs or changing URLs that merely share a prefix.
func ReplaceURLs(content, oldURL, newURL string) (string, int) {
	if oldURL == "" || oldURL == newURL {
		return content, 0
	}
	var out strings.Builder
	count := 0
	for len(content) > 0 {
		i := strings.Index(content, oldURL)
		if i < 0 {
			out.WriteString(content)
			break
		}
		end := i + len(oldURL)
		beforeOK := i == 0 || strings.ContainsRune(" \t\r\n(<\"'=", rune(content[i-1]))
		afterOK := end == len(content) || strings.ContainsRune(" \t\r\n)>\"'", rune(content[end]))
		// Bare URLs may be followed by sentence punctuation.
		if end < len(content) && strings.ContainsRune(".,;!", rune(content[end])) {
			afterOK = end+1 == len(content) || unicode.IsSpace(rune(content[end+1]))
		}
		out.WriteString(content[:i])
		if beforeOK && afterOK {
			out.WriteString(newURL)
			count++
		} else {
			out.WriteString(oldURL)
		}
		content = content[end:]
	}
	return out.String(), count
}
