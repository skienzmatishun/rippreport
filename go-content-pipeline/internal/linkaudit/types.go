// Package linkaudit implements the audit/link_audit workflows without the
// original Python scripts. Only optional archive rendering uses Playwright.
package linkaudit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/archive"
)

type Link struct {
	URL          string  `json:"url"`
	Text         string  `json:"text"`
	OriginalHref string  `json:"original_href"`
	Status       string  `json:"status"`
	StatusCode   *int    `json:"status_code"`
	FinalURL     *string `json:"final_url"`
	Redirects    int     `json:"redirects"`
	Error        *string `json:"error"`
	Broken       bool    `json:"is_broken"`
}

type Post struct {
	URL         string  `json:"post_url"`
	Accessible  bool    `json:"post_accessible"`
	Error       *string `json:"error"`
	Links       []Link  `json:"links"`
	TotalLinks  int     `json:"total_links"`
	BrokenLinks int     `json:"broken_links"`
}

type Summary struct {
	TotalPosts        int `json:"total_posts"`
	TotalLinks        int `json:"total_links"`
	BrokenLinks       int `json:"broken_links"`
	InaccessiblePosts int `json:"inaccessible_posts"`
}

type Report struct {
	Posts   []Post  `json:"posts"`
	Summary Summary `json:"summary"`
}

func (r *Report) Recount() {
	r.Summary = Summary{TotalPosts: len(r.Posts)}
	for i := range r.Posts {
		p := &r.Posts[i]
		p.TotalLinks = len(p.Links)
		p.BrokenLinks = 0
		if !p.Accessible {
			r.Summary.InaccessiblePosts++
		}
		for _, l := range p.Links {
			if l.Broken {
				p.BrokenLinks++
			}
		}
		r.Summary.TotalLinks += p.TotalLinks
		r.Summary.BrokenLinks += p.BrokenLinks
	}
}

func LoadReport(path string) (*Report, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var shape map[string]json.RawMessage
	if err := json.Unmarshal(data, &shape); err != nil {
		return nil, err
	}
	raw, ok := shape["posts"]
	if !ok || string(raw) == "null" {
		return nil, fmt.Errorf("link report must contain a posts array")
	}
	var report Report
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}
	for _, p := range report.Posts {
		if p.URL == "" {
			return nil, fmt.Errorf("link report contains a post without post_url")
		}
	}
	report.Recount()
	return &report, nil
}

func webURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Host != "" && (u.Scheme == "http" || u.Scheme == "https")
}
func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
func message(s string) *string { return &s }

func bundles(root string, slugs []string) ([]string, error) {
	items, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	selected := map[string]bool{}
	for _, slug := range slugs {
		selected[slug] = true
	}
	var result []string
	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		if len(slugs) > 0 && !selected[item.Name()] {
			continue
		}
		raw := "/p/" + url.PathEscape(item.Name()) + "/"
		if _, err := archive.PostPath(root, raw); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("bundle %s: %w", item.Name(), err)
		}
		result = append(result, item.Name())
		delete(selected, item.Name())
	}
	if len(selected) > 0 {
		var missing []string
		for s := range selected {
			missing = append(missing, s)
		}
		sort.Strings(missing)
		return nil, fmt.Errorf("post bundles not found: %s", strings.Join(missing, ", "))
	}
	sort.Strings(result)
	return result, nil
}

func distinct(paths ...string) error {
	seen := map[string]bool{}
	for _, p := range paths {
		if p == "" {
			continue
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			return err
		}
		if resolved, err := filepath.EvalSymlinks(abs); err == nil {
			abs = resolved
		}
		if seen[abs] {
			return fmt.Errorf("input and output paths must be distinct: %s", p)
		}
		seen[abs] = true
	}
	return nil
}
