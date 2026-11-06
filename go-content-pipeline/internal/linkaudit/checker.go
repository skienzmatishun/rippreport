package linkaudit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/archive"
	"golang.org/x/net/html"
)

type CheckOptions struct {
	ContentDir, BaseURL, Report            string
	Posts                                  []string
	SkipDomains, SkipPaths                 []string
	Timeout, Delay                         time.Duration
	MaxRedirects, PostWorkers, LinkWorkers int
	Resume                                 bool
}

func DefaultCheckOptions() CheckOptions {
	return CheckOptions{ContentDir: "content/p", BaseURL: "http://localhost:1313", Report: "audit/broken_links_report.json", Timeout: 10 * time.Second, Delay: 500 * time.Millisecond, MaxRedirects: 10, PostWorkers: 5, LinkWorkers: 20, Resume: true,
		SkipDomains: []string{"google.com", "facebook.com", "fb.com", "twitter.com", "x.com", "localhost", "sibforms.com"}, SkipPaths: []string{"/p/privacy-policy/", "/privacy-policy/"}}
}

func ShouldSkip(raw string, domains, paths []string) bool {
	u, err := url.Parse(raw)
	if err != nil || !webURL(raw) {
		return true
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	for _, domain := range domains {
		domain = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(domain)), "www.")
		if domain != "" && (host == domain || strings.HasSuffix(host, "."+domain)) {
			return true
		}
	}
	for _, p := range paths {
		if p != "" && strings.HasSuffix(strings.ToLower(u.Path), strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// ExtractLinks resolves all relative destinations and retains original hrefs for
// exact replacement in source Markdown. Non-HTTP schemes are never fetched.
func ExtractLinks(source, base string, domains, paths []string) ([]Link, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	root, err := html.Parse(strings.NewReader(source))
	if err != nil {
		return nil, err
	}
	links := []Link{}
	var visit func(*html.Node)
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				href := strings.TrimSpace(a.Val)
				if href == "" || strings.HasPrefix(href, "#") {
					continue
				}
				rel, err := url.Parse(href)
				if err != nil {
					continue
				}
				absolute := u.ResolveReference(rel).String()
				if ShouldSkip(absolute, domains, paths) {
					continue
				}
				var text strings.Builder
				var collect func(*html.Node)
				collect = func(c *html.Node) {
					if c.Type == html.TextNode {
						text.WriteString(c.Data)
					}
					for child := c.FirstChild; child != nil; child = child.NextSibling {
						collect(child)
					}
				}
				collect(n)
				label := []rune(strings.TrimSpace(text.String()))
				if len(label) > 100 {
					label = label[:100]
				}
				links = append(links, Link{URL: absolute, Text: string(label), OriginalHref: a.Val})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(root)
	return links, nil
}

var errRedirects = errors.New("too many redirects")

type Checker struct {
	Options CheckOptions
	HTTP    *http.Client
	Log     io.Writer
	slots   chan struct{}
}

func NewChecker(o CheckOptions, log io.Writer) *Checker {
	if log == nil {
		log = io.Discard
	}
	return &Checker{Options: o, HTTP: &http.Client{}, Log: log}
}

func (c *Checker) fetch(ctx context.Context, raw string, readBody bool) ([]byte, Link) {
	result := Link{URL: raw, Status: "ok"}
	ctx, cancel := context.WithTimeout(ctx, c.Options.Timeout)
	defer cancel()
	client := *c.HTTP
	redirects := 0
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > c.Options.MaxRedirects {
			return errRedirects
		}
		if !webURL(req.URL.String()) {
			return fmt.Errorf("non-HTTP redirect")
		}
		redirects++
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, "GET", raw, nil)
	if err == nil {
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; LinkChecker/2.0)")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	}
	var resp *http.Response
	if err == nil {
		resp, err = client.Do(req)
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	result.Redirects = redirects
	if err != nil {
		result.Broken = true
		result.Status = "connection_error"
		result.Error = message(err.Error())
		var netErr net.Error
		if errors.Is(err, errRedirects) {
			result.Status = "too_many_redirects"
		} else if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout()) {
			result.Status = "timeout"
		} else if errors.Is(err, context.Canceled) {
			result.Status = "cancelled"
		}
		return nil, result
	}
	result.StatusCode = &resp.StatusCode
	final := resp.Request.URL.String()
	result.FinalURL = &final
	result.Broken = resp.StatusCode >= 400
	if !readBody {
		return nil, result
	}
	const limit = 16 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil || len(body) > limit {
		result.Status = "error"
		result.Broken = true
		if err != nil {
			result.Error = message(err.Error())
		} else {
			result.Error = message("post exceeds 16 MiB")
		}
		return nil, result
	}
	return body, result
}

func (c *Checker) checkPost(ctx context.Context, raw string) Post {
	p := Post{URL: raw, Links: []Link{}}
	body, response := c.fetch(ctx, raw, true)
	if response.StatusCode == nil || *response.StatusCode != 200 || response.Error != nil {
		p.Error = response.Error
		if p.Error == nil {
			p.Error = message(fmt.Sprintf("Could not fetch post: %d", *response.StatusCode))
		}
		return p
	}
	base := raw
	if response.FinalURL != nil {
		base = *response.FinalURL
	}
	links, err := ExtractLinks(string(body), base, c.Options.SkipDomains, c.Options.SkipPaths)
	if err != nil {
		p.Error = message(err.Error())
		return p
	}
	p.Accessible = true
	jobs := make(chan int)
	var wg sync.WaitGroup
	for worker := 0; worker < c.Options.LinkWorkers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				select {
				case c.slots <- struct{}{}:
				case <-ctx.Done():
					return
				}
				_, result := c.fetch(ctx, links[i].URL, false)
				result.Text = links[i].Text
				result.OriginalHref = links[i].OriginalHref
				links[i] = result
				_ = wait(ctx, c.Options.Delay)
				<-c.slots
			}
		}()
	}
send:
	for i := range links {
		select {
		case jobs <- i:
		case <-ctx.Done():
			break send
		}
	}
	close(jobs)
	wg.Wait()
	p.Links = links
	for _, l := range links {
		if l.Broken {
			p.BrokenLinks++
		}
	}
	p.TotalLinks = len(links)
	return p
}

func (c *Checker) Run(ctx context.Context) (*Report, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	o := c.Options
	if !webURL(o.BaseURL) || o.Report == "" || o.ContentDir == "" || o.Timeout <= 0 || o.Delay < 0 || o.MaxRedirects < 0 || o.PostWorkers < 1 || o.LinkWorkers < 1 {
		return nil, fmt.Errorf("invalid link checker options")
	}
	names, err := bundles(o.ContentDir, o.Posts)
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("no post bundles found in %s", o.ContentDir)
	}
	unlock, err := archive.Lock(ctx, o.Report+".lock", 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer unlock()
	report := &Report{Posts: []Post{}}
	if loaded, err := LoadReport(o.Report); err == nil {
		report = loaded
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	positions := map[string]int{}
	for i, p := range report.Posts {
		positions[p.URL] = i
	}
	var scheduled []string
	for _, name := range names {
		raw := strings.TrimRight(o.BaseURL, "/") + "/p/" + url.PathEscape(name) + "/"
		if _, done := positions[raw]; !done || !o.Resume {
			scheduled = append(scheduled, raw)
		}
	}
	if len(scheduled) == 0 {
		return report, nil
	}
	c.slots = make(chan struct{}, o.LinkWorkers)
	jobs := make(chan string)
	results := make(chan Post)
	var wg sync.WaitGroup
	for worker := 0; worker < o.PostWorkers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for raw := range jobs {
				if ctx.Err() != nil {
					return
				}
				p := c.checkPost(ctx, raw)
				if ctx.Err() == nil {
					results <- p
				}
			}
		}()
	}
	go func() {
		defer close(jobs)
		for _, raw := range scheduled {
			select {
			case jobs <- raw:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() { wg.Wait(); close(results) }()
	var saveErr error
	for p := range results {
		if i, ok := positions[p.URL]; ok {
			report.Posts[i] = p
		} else {
			positions[p.URL] = len(report.Posts)
			report.Posts = append(report.Posts, p)
		}
		report.Recount()
		fmt.Fprintf(c.Log, "%s: %d links, %d broken, accessible=%t\n", p.URL, p.TotalLinks, p.BrokenLinks, p.Accessible)
		if saveErr == nil {
			saveErr = archive.SaveJSON(o.Report, report)
		}
	}
	sort.Slice(report.Posts, func(i, j int) bool { return report.Posts[i].URL < report.Posts[j].URL })
	report.Recount()
	if saveErr != nil {
		return report, saveErr
	}
	if err := archive.SaveJSON(o.Report, report); err != nil {
		return report, err
	}
	if ctx.Err() != nil {
		return report, ctx.Err()
	}
	return report, nil
}
