package archive

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

//go:embed browser.py
var browserScript string

type page struct {
	HTML               string `json:"html"`
	URL                string `json:"url"`
	Status             int    `json:"status"`
	ScreenshotError    string `json:"screenshot_error"`
	Screenshot, Method string
}

type Verifier struct {
	Options    Options
	HTTP       *http.Client
	Log        io.Writer
	lastSearch time.Time
}

func New(o Options, log io.Writer) *Verifier {
	if log == nil {
		log = io.Discard
	}
	return &Verifier{Options: o, HTTP: &http.Client{Timeout: o.RequestTimeout}, Log: log}
}

func pause(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (v *Verifier) request(ctx context.Context, method, raw string) ([]byte, http.Header, string, error) {
	if !webURL(raw) {
		return nil, nil, "", fmt.Errorf("invalid HTTP(S) URL: %s", raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, raw, nil)
	if err != nil {
		return nil, nil, "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ArchiveVerifier/2.0)")
	resp, err := v.HTTP.Do(req)
	if err != nil {
		return nil, nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.Header, resp.Request.URL.String(), fmt.Errorf("HTTP %d for %s", resp.StatusCode, raw)
	}
	const maxBody = 16 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if len(body) > maxBody {
		return nil, resp.Header, "", fmt.Errorf("page exceeds 16 MiB")
	}
	return body, resp.Header, resp.Request.URL.String(), err
}

func (v *Verifier) load(ctx context.Context, raw string) (page, error) {
	if !webURL(raw) {
		return page{}, fmt.Errorf("invalid HTTP(S) URL: %s", raw)
	}
	if v.Options.UsePlaywright {
		p, err := v.render(ctx, raw)
		if err == nil {
			return p, nil
		}
		fmt.Fprintf(v.Log, "Browser unavailable for %s: %v; using HTTP\n", raw, err)
	}
	b, _, final, err := v.request(ctx, "GET", raw)
	return page{HTML: string(b), URL: final, Status: 200, Method: "http"}, err
}

func (v *Verifier) render(ctx context.Context, raw string) (page, error) {
	o := v.Options
	var screenshot string
	if o.UseLMStudio && o.Vision {
		if err := os.MkdirAll(o.ScreenshotDir, 0755); err != nil {
			return page{}, err
		}
		f, err := os.CreateTemp(o.ScreenshotDir, "archive-*.png")
		if err != nil {
			return page{}, err
		}
		screenshot, _ = filepath.Abs(f.Name())
		f.Close()
	}
	ok := false
	defer func() {
		if !ok && screenshot != "" {
			os.Remove(screenshot)
		}
	}()
	body, _ := json.Marshal(map[string]any{"url": raw, "timeout_ms": o.BrowserTimeout.Milliseconds(), "wait_ms": o.Wait.Milliseconds(), "screenshot": screenshot})
	ctx, cancel := context.WithTimeout(ctx, o.BrowserTimeout+o.Wait+15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, o.Python, "-c", browserScript)
	cmd.Stdin = bytes.NewReader(body)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return page{}, fmt.Errorf("Playwright: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	var p page
	if err := json.Unmarshal(output, &p); err != nil {
		return page{}, fmt.Errorf("browser response: %w", err)
	}
	if p.Status >= 400 {
		return page{}, fmt.Errorf("browser HTTP %d", p.Status)
	}
	p.Method = "playwright"
	if p.ScreenshotError == "" {
		p.Screenshot = screenshot
	} else {
		os.Remove(screenshot)
	}
	ok = true
	return p, nil
}

func (v *Verifier) clean(p page) {
	if v.Options.CleanupScreenshots && p.Screenshot != "" {
		os.Remove(p.Screenshot)
	}
}

// analyze uses the OpenAI-compatible chat API supported by LM Studio, Unsloth,
// and llama-server. Invalid or unavailable model responses fall back to text.
func (v *Verifier) analyze(ctx context.Context, prompt, image string) (map[string]any, error) {
	return v.analyzeJSON(ctx, prompt, image, "is_valid")
}

// AnalyzeContentJSON shares model transport, thinking/fence handling, retries,
// and context cancellation with archive analysis, using a content-audit schema.
func (v *Verifier) AnalyzeContentJSON(ctx context.Context, prompt string) (map[string]any, error) {
	return v.analyzeJSON(ctx, prompt, "", "has_issue")
}

func (v *Verifier) analyzeJSON(ctx context.Context, prompt, image, booleanField string) (map[string]any, error) {
	o := v.Options
	var content any = prompt
	if image != "" && o.Vision {
		b, err := os.ReadFile(image)
		if err != nil {
			return nil, err
		}
		content = []any{map[string]any{"type": "text", "text": prompt}, map[string]any{"type": "image_url", "image_url": map[string]string{"url": "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)}}}
	}
	payload, _ := json.Marshal(map[string]any{"model": o.Model, "temperature": 0, "max_tokens": 600, "messages": []any{map[string]any{"role": "user", "content": content}}})
	endpoint := strings.TrimRight(o.ModelURL, "/")
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		if !strings.HasSuffix(endpoint, "/v1") {
			endpoint += "/v1"
		}
		endpoint += "/chat/completions"
	}
	var last error
	for attempt := 0; attempt < o.Retries; attempt++ {
		if attempt > 0 {
			if err := pause(ctx, time.Second*time.Duration(1<<uint(attempt-1))); err != nil {
				return nil, err
			}
		}
		callCtx, cancel := context.WithTimeout(ctx, o.ModelTimeout)
		req, err := http.NewRequestWithContext(callCtx, "POST", endpoint, bytes.NewReader(payload))
		if err != nil {
			cancel()
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer not-needed")
		// Use a separate timeout: page-fetch timeout must not shorten model calls.
		client := &http.Client{Transport: v.HTTP.Transport, Timeout: o.ModelTimeout}
		resp, err := client.Do(req)
		if err != nil {
			cancel()
			last = err
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			continue
		}
		b, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		resp.Body.Close()
		cancel()
		if readErr != nil {
			last = readErr
			continue
		}
		if resp.StatusCode != 200 {
			last = fmt.Errorf("model HTTP %d", resp.StatusCode)
			if resp.StatusCode < 500 && resp.StatusCode != 429 {
				break
			}
			continue
		}
		var response struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.Unmarshal(b, &response); err != nil {
			last = err
			continue
		}
		if len(response.Choices) == 0 {
			last = fmt.Errorf("empty model response")
			continue
		}
		s := response.Choices[0].Message.Content
		if end := strings.LastIndex(s, "</think>"); end >= 0 {
			s = s[end+len("</think>"):]
		}
		first, end := strings.Index(s, "{"), strings.LastIndex(s, "}")
		if first < 0 || end < first {
			last = fmt.Errorf("model returned no JSON object")
			continue
		}
		var result map[string]any
		if err := json.Unmarshal([]byte(s[first:end+1]), &result); err != nil {
			last = err
			continue
		}
		if _, ok := result[booleanField].(bool); !ok {
			last = fmt.Errorf("model response missing boolean %s", booleanField)
			continue
		}
		return result, nil
	}
	return nil, last
}

func (v *Verifier) verify(ctx context.Context, raw string, e Entry, paywall bool) Verification {
	r := Verification{Method: "failed", VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
	p, err := v.load(ctx, raw)
	if err != nil {
		r.Reason = err.Error()
		return r
	}
	defer v.clean(p)
	d := inspectHTML(p.HTML)
	r.Method = p.Method + "+html"
	r.Valid = d.substantial(paywall)
	r.Reason = "No substantial article content detected"
	if r.Valid {
		r.Reason = "Substantial article content found"
	}
	if !v.Options.UseLMStudio {
		return r
	}
	contextJSON, _ := json.Marshal(map[string]string{"original_url": e.str("original_href"), "url": raw, "title": e.str("post_title"), "date": e.str("post_date"), "excerpt": e.str("post_excerpt")})
	text := []rune(d.Text)
	if len(text) > v.Options.MaxContent {
		text = text[:v.Options.MaxContent]
	}
	prompt := "Verify this page contains the actual requested article rather than a blank page, error, homepage, bot challenge, or subscription gate. Check relevance to the supplied context. Return only JSON: {\"is_valid\":boolean,\"confidence\":number,\"reason\":string}. Context: " + string(contextJSON) + "\nPage text:\n" + string(text)
	analysis, err := v.analyze(ctx, prompt, p.Screenshot)
	usedVision := p.Screenshot != ""
	if err != nil && p.Screenshot != "" {
		usedVision = false
		analysis, err = v.analyze(ctx, prompt, "")
	}
	if err != nil {
		fmt.Fprintf(v.Log, "Model analysis unavailable: %v; using HTML analysis\n", err)
		return r
	}
	r.Valid, _ = analysis["is_valid"].(bool)
	r.Reason, _ = analysis["reason"].(string)
	r.Confidence, _ = analysis["confidence"].(float64)
	r.Method = p.Method + "+lmstudio"
	if usedVision {
		r.Method += "+vision"
	}
	r.Details = analysis
	return r
}

var timestampPattern = regexp.MustCompile(`/web/(\d{14})(?:[a-z_]+)?/`)

func (v *Verifier) snapshots(ctx context.Context, original, before string) ([]string, error) {
	u, err := url.Parse(v.Options.CDXURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("url", original)
	q.Set("output", "json")
	q.Set("fl", "timestamp,original,statuscode")
	q.Set("filter", "statuscode:200")
	q.Set("to", before)
	q.Set("limit", fmt.Sprint(-v.Options.MaxSnapshots))
	u.RawQuery = q.Encode()
	b, _, _, err := v.request(ctx, "GET", u.String())
	if err != nil {
		return nil, err
	}
	var rows [][]string
	if err := json.Unmarshal(b, &rows); err != nil {
		return nil, err
	}
	var valid [][]string
	for _, row := range rows {
		if len(row) == 3 && row[0] < before && len(row[0]) == 14 && row[2] == "200" && webURL(row[1]) {
			if _, err := time.Parse("20060102150405", row[0]); err != nil {
				continue
			}
			valid = append(valid, row)
		}
	}
	sort.Slice(valid, func(i, j int) bool { return valid[i][0] > valid[j][0] })
	var result []string
	for _, row := range valid {
		if len(result) == v.Options.MaxSnapshots {
			break
		}
		result = append(result, "https://web.archive.org/web/"+row[0]+"/"+row[1])
	}
	return result, nil
}

func (v *Verifier) search(ctx context.Context, query string) ([]string, error) {
	if !v.Options.Google {
		return nil, fmt.Errorf("Google fallback disabled")
	}
	if !v.lastSearch.IsZero() {
		jitter := time.Duration(0)
		if v.Options.SearchJitter > 0 {
			jitter = time.Duration(rand.Int64N(int64(v.Options.SearchJitter)))
		}
		if err := pause(ctx, time.Until(v.lastSearch.Add(v.Options.SearchDelay+jitter))); err != nil {
			return nil, err
		}
	}
	v.lastSearch = time.Now()
	u, _ := url.Parse(v.Options.SearchURL)
	q := u.Query()
	q.Set("q", query)
	q.Set("num", fmt.Sprint(v.Options.MaxResults))
	u.RawQuery = q.Encode()
	b, _, _, err := v.request(ctx, "GET", u.String())
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var result []string
	for _, href := range inspectHTML(string(b)).Links {
		link, err := url.Parse(href)
		if err != nil {
			continue
		}
		if link.Path == "/url" {
			href = link.Query().Get("q")
			if href == "" {
				href = link.Query().Get("url")
			}
			link, err = url.Parse(href)
		}
		if err != nil || !webURL(href) || link.Hostname() == "google.com" || strings.HasSuffix(link.Hostname(), ".google.com") || seen[href] {
			continue
		}
		seen[href] = true
		result = append(result, href)
		if len(result) >= v.Options.MaxResults {
			break
		}
	}
	return result, nil
}

func (v *Verifier) verifyPDF(ctx context.Context, raw string) Verification {
	r := Verification{Method: "content-type", VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
	_, headers, _, err := v.request(ctx, "HEAD", raw)
	if err != nil || !strings.Contains(strings.ToLower(headers.Get("Content-Type")), "application/pdf") {
		r.Reason = "URL does not return a successful PDF content type"
		return r
	}
	r.Valid = true
	r.Reason = "Valid PDF content-type header"
	if v.Options.UsePlaywright && v.Options.UseLMStudio && v.Options.Vision {
		p, err := v.render(ctx, raw)
		if err == nil {
			defer v.clean(p)
			if p.Screenshot != "" {
				a, err := v.analyze(ctx, "Does this screenshot show a readable PDF document rather than an error or blank page? Return JSON {\"is_valid\":boolean,\"confidence\":number,\"reason\":string}. PDF URL: "+raw, p.Screenshot)
				if err == nil {
					r.Valid, _ = a["is_valid"].(bool)
					r.Reason, _ = a["reason"].(string)
					r.Confidence, _ = a["confidence"].(float64)
					r.Method = "vision"
					r.Details = a
				}
			}
		}
	}
	return r
}

// Process selects and verifies a replacement; it never modifies content.
func (v *Verifier) Process(ctx context.Context, e Entry) (string, Verification) {
	failure := func(reason string) (string, Verification) {
		return "", Verification{Reason: reason, Method: "failed", VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
	}
	original := e.str("original_href")
	switch Strategy(e, v.Options) {
	case "verify_fixed":
		raw := e.str("archive_url")
		r := v.verify(ctx, raw, e, false)
		if r.Valid {
			return raw, r
		}
		m := timestampPattern.FindStringSubmatch(raw)
		if m == nil {
			return failure("Invalid archive and no Wayback timestamp found")
		}
		snapshots, err := v.snapshots(ctx, original, m[1])
		if err != nil {
			return failure("Wayback lookup: " + err.Error())
		}
		for _, candidate := range snapshots {
			if ctx.Err() != nil {
				break
			}
			r = v.verify(ctx, candidate, e, false)
			if r.Valid {
				return candidate, r
			}
		}
		return failure("All archive snapshots invalid")
	case "pdf_fallback":
		results, err := v.search(ctx, original)
		if err != nil {
			return failure(err.Error())
		}
		for _, raw := range results {
			r := v.verifyPDF(ctx, raw)
			if r.Valid {
				return raw, r
			}
		}
		return failure("No valid relocated PDF found")
	case "malformed_cdn":
		raw := strings.TrimSpace(e.str("link_text"))
		if raw == "" {
			raw = strings.TrimSpace(e.str("text"))
		}
		if webURL(raw) {
			if _, _, _, err := v.request(ctx, "HEAD", raw); err == nil {
				return raw, Verification{Valid: true, Method: "link_text", Reason: "Link text URL is accessible", VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
			}
		}
		if v.Options.UseLMStudio {
			b, _ := json.Marshal(e)
			a, err := v.analyze(ctx, "Recover the intended URL from this malformed CDN link and context. Do not invent a URL. Return JSON {\"is_valid\":boolean,\"new_url\":string,\"reason\":string}: "+string(b), "")
			if err == nil {
				candidate, _ := a["new_url"].(string)
				ok, _ := a["is_valid"].(bool)
				if ok && webURL(candidate) {
					r := v.verify(ctx, candidate, e, false)
					if r.Valid {
						return candidate, r
					}
				}
			}
		}
		return failure("No verified replacement for malformed CDN link")
	case "tag_conversion":
		u, _ := url.Parse(original)
		target := v.Options.SearchURL + "?q=" + url.QueryEscape("site:rippreport.com "+tagName(u.Path))
		return target, Verification{Valid: true, Method: "tag_conversion", Reason: "Converted tag to Google site search", VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
	case "missing_p_path":
		u, _ := url.Parse(original)
		slug := missingSlug(u.Path)
		// A local bundle proves the canonical URL without needing a search.
		candidate := "https://rippreport.com/p/" + url.PathEscape(slug) + "/"
		if _, err := PostPath(v.Options.ContentDir, candidate); err == nil {
			return candidate, Verification{Valid: true, Method: "local_bundle", Reason: "Canonical post bundle exists", VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
		}
		results, err := v.search(ctx, "site:rippreport.com "+slug)
		if err != nil {
			return failure(err.Error())
		}
		for _, raw := range results {
			u, err := url.Parse(raw)
			if err != nil || !siteHost(u.Hostname()) || !strings.HasPrefix(u.Path, "/p/") {
				continue
			}
			r := Verification{Valid: true, Method: "google", Reason: "Found canonical post URL", VerifiedAt: time.Now().UTC().Format(time.RFC3339)}
			if v.Options.VerifyGoogle {
				r = v.verify(ctx, raw, e, false)
			}
			if r.Valid {
				return raw, r
			}
		}
		return failure("No verified canonical post URL found")
	case "paywall_check":
		r := v.verify(ctx, original, e, true)
		if r.Valid {
			return original, r
		}
		return "", r
	}
	return failure("No applicable strategy")
}
