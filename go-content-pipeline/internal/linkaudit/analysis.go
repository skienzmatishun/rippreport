package linkaudit

import (
	"fmt"
	"io"
	"net/url"
	"sort"
)

type Reference struct {
	Post string `json:"post"`
	Link Link   `json:"link"`
}
type Analysis struct {
	Summary     Summary                `json:"summary"`
	SuccessRate float64                `json:"success_rate"`
	ByDomain    map[string][]Reference `json:"broken_by_domain"`
	ByStatus    map[string][]Reference `json:"broken_by_status"`
	ByError     map[string][]Reference `json:"broken_by_error"`
	TopPosts    []Post                 `json:"posts_with_most_broken_links"`
}

func Analyze(report *Report) Analysis {
	report.Recount()
	a := Analysis{Summary: report.Summary, ByDomain: map[string][]Reference{}, ByStatus: map[string][]Reference{}, ByError: map[string][]Reference{}, TopPosts: []Post{}}
	if a.Summary.TotalLinks > 0 {
		a.SuccessRate = float64(a.Summary.TotalLinks-a.Summary.BrokenLinks) / float64(a.Summary.TotalLinks) * 100
	}
	for _, p := range report.Posts {
		if p.BrokenLinks > 0 {
			a.TopPosts = append(a.TopPosts, p)
		}
		for _, l := range p.Links {
			if !l.Broken {
				continue
			}
			domain := "unknown"
			if u, err := url.Parse(l.URL); err == nil && u.Host != "" {
				domain = u.Host
			}
			status := "unknown"
			if l.StatusCode != nil {
				status = fmt.Sprint(*l.StatusCode)
			}
			errorType := l.Status
			if errorType == "" {
				errorType = "unknown"
			}
			ref := Reference{Post: p.URL, Link: l}
			a.ByDomain[domain] = append(a.ByDomain[domain], ref)
			a.ByStatus[status] = append(a.ByStatus[status], ref)
			a.ByError[errorType] = append(a.ByError[errorType], ref)
		}
	}
	sort.Slice(a.TopPosts, func(i, j int) bool {
		if a.TopPosts[i].BrokenLinks == a.TopPosts[j].BrokenLinks {
			return a.TopPosts[i].URL < a.TopPosts[j].URL
		}
		return a.TopPosts[i].BrokenLinks > a.TopPosts[j].BrokenLinks
	})
	if len(a.TopPosts) > 10 {
		a.TopPosts = a.TopPosts[:10]
	}
	return a
}

func (a Analysis) Print(w io.Writer) {
	fmt.Fprintf(w, "Posts checked: %d (%d inaccessible)\nTotal links: %d\nBroken links: %d\nSuccess rate: %.1f%%\n", a.Summary.TotalPosts, a.Summary.InaccessiblePosts, a.Summary.TotalLinks, a.Summary.BrokenLinks, a.SuccessRate)
	for _, group := range []struct {
		name  string
		items map[string][]Reference
		limit int
	}{{"domain", a.ByDomain, 5}, {"status code", a.ByStatus, 3}, {"error type", a.ByError, 3}} {
		var keys []string
		for k := range group.items {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			if group.name == "domain" && len(group.items[keys[i]]) != len(group.items[keys[j]]) {
				return len(group.items[keys[i]]) > len(group.items[keys[j]])
			}
			return keys[i] < keys[j]
		})
		fmt.Fprintf(w, "\nBroken links by %s:\n", group.name)
		for _, k := range keys {
			refs := group.items[k]
			fmt.Fprintf(w, "  %s: %d\n", k, len(refs))
			for i, ref := range refs {
				if i == group.limit {
					fmt.Fprintf(w, "    ... and %d more\n", len(refs)-i)
					break
				}
				fmt.Fprintf(w, "    %s (post: %s)\n", ref.Link.URL, ref.Post)
				if ref.Link.Error != nil {
					fmt.Fprintf(w, "    Error: %s\n", *ref.Link.Error)
				}
			}
		}
	}
	fmt.Fprintln(w, "\nPosts with most broken links:")
	for _, p := range a.TopPosts {
		fmt.Fprintf(w, "  %s: %d/%d broken\n", p.URL, p.BrokenLinks, p.TotalLinks)
	}
}
