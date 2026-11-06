package cli

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/archive"
	"github.com/rippreport/go-content-pipeline/internal/linkaudit"
)

func auditPaths() (archive.Options, string) {
	o := archiveDefaults()
	return o, filepath.Join(filepath.Dir(o.Report), "broken_links_report.json")
}

func (a *App) auditError(ctx context.Context, err error) int {
	if err == nil {
		return 0
	}
	fmt.Fprintf(a.stderr, "Audit failed: %v\n", err)
	if ctx.Err() != nil {
		return 130
	}
	return 1
}

func (a *App) runCheckLinks(ctx context.Context, args []string) int {
	defaults, input := auditPaths()
	o := linkaudit.DefaultCheckOptions()
	o.ContentDir = defaults.ContentDir
	o.Report = input
	fs := flag.NewFlagSet("check-links", flag.ContinueOnError)
	fs.SetOutput(a.stdout)
	fs.StringVar(&o.ContentDir, "content-dir", o.ContentDir, "Hugo post bundle directory")
	fs.StringVar(&o.BaseURL, "base-url", o.BaseURL, "Running Hugo site URL")
	fs.StringVar(&o.Report, "report", o.Report, "Resumable link audit JSON")
	fs.DurationVar(&o.Timeout, "request-timeout", o.Timeout, "HTTP timeout including redirects")
	fs.DurationVar(&o.Delay, "request-delay", o.Delay, "Delay after each link request per worker")
	fs.IntVar(&o.MaxRedirects, "max-redirects", o.MaxRedirects, "Maximum redirects per URL")
	fs.IntVar(&o.PostWorkers, "post-workers", o.PostWorkers, "Concurrent post workers")
	fs.IntVar(&o.LinkWorkers, "link-workers", o.LinkWorkers, "Maximum concurrent link requests across all posts")
	fs.BoolVar(&o.Resume, "resume", o.Resume, "Skip posts already in report; false rechecks selected posts")
	posts := fs.String("posts", "", "Comma-separated post slugs")
	domains := fs.String("skip-domains", strings.Join(o.SkipDomains, ","), "Comma-separated excluded domains and their subdomains; empty disables")
	paths := fs.String("skip-paths", strings.Join(o.SkipPaths, ","), "Comma-separated excluded path suffixes; empty disables")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(a.stderr, "Unexpected positional arguments")
		return 1
	}
	o.Posts = parseSlugs(*posts)
	o.SkipDomains = parseSlugs(*domains)
	o.SkipPaths = parseSlugs(*paths)
	report, err := linkaudit.NewChecker(o, a.stdout).Run(ctx)
	if report != nil {
		linkaudit.Analyze(report).Print(a.stdout)
	}
	return a.auditError(ctx, err)
}

func (a *App) runAnalyzeBrokenLinks(ctx context.Context, args []string) int {
	_, input := auditPaths()
	fs := flag.NewFlagSet("analyze-broken-links", flag.ContinueOnError)
	fs.SetOutput(a.stdout)
	reportPath := fs.String("report", input, "Link audit JSON")
	output := fs.String("output", "", "Optional JSON with groups by domain, status and error")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(a.stderr, "Unexpected positional arguments")
		return 1
	}
	if *output != "" && sameFilePath(*reportPath, *output) {
		fmt.Fprintln(a.stderr, "Output must differ from input report")
		return 1
	}
	report, err := linkaudit.LoadReport(*reportPath)
	if err != nil {
		return a.auditError(ctx, err)
	}
	analysis := linkaudit.Analyze(report)
	analysis.Print(a.stdout)
	if *output != "" {
		return a.auditError(ctx, archive.SaveJSON(*output, analysis))
	}
	return 0
}

func (a *App) runLinkRepairs(ctx context.Context, args []string, redirects bool) int {
	defaults, input := auditPaths()
	command, filename := "archive-links", "archive_fix_report.json"
	if redirects {
		command, filename = "update-redirected-links", "redirected_links_update_report.json"
	}
	o := linkaudit.RepairOptions{Input: input, Output: filepath.Join(filepath.Dir(input), filename), Archive: defaults, MinRedirects: 1, Resume: true, CaseInsensitive: true}
	o.Archive.UsePlaywright = false
	o.Archive.UseLMStudio = false
	o.Archive.Google = false
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(a.stdout)
	fs.StringVar(&o.Input, "report", o.Input, "Source link audit JSON")
	fs.StringVar(&o.Output, "output", o.Output, "Repair report JSON")
	fs.StringVar(&o.Archive.ContentDir, "content-dir", o.Archive.ContentDir, "Hugo post bundle directory")
	fs.StringVar(&o.Archive.BackupDir, "backup-dir", o.Archive.BackupDir, "Verified backups before content changes")
	dryRun := fs.Bool("dry-run", false, "Preview without writing content, backups or output reports")
	var update *bool
	if redirects {
		update = fs.Bool("update", false, "Apply changes (default previews only)")
		fs.IntVar(&o.MinRedirects, "min-redirects", 1, "Minimum redirect count")
		fs.StringVar(&o.Domain, "domain", "", "Filter original domain by case-insensitive substring")
		fs.BoolVar(&o.CaseInsensitive, "case-insensitive", true, "Allow exact URL destinations with case differences")
	} else {
		o.CSV = filepath.Join(filepath.Dir(input), "archive_fix_report.csv")
		fs.StringVar(&o.CSV, "csv", o.CSV, "CSV repair report; empty disables")
		fs.StringVar(&o.Archive.CDXURL, "cdx-url", o.Archive.CDXURL, "Wayback CDX endpoint")
		fs.IntVar(&o.Archive.MaxSnapshots, "max-snapshots", o.Archive.MaxSnapshots, "Maximum snapshots to try, newest first")
		fs.BoolVar(&o.Resume, "resume", true, "Skip previously attempted links; false retries")
		fs.BoolVar(&o.Archive.UsePlaywright, "playwright", false, "Render snapshots using Playwright")
		fs.StringVar(&o.Archive.Python, "python", o.Archive.Python, "Python interpreter with Playwright")
		fs.BoolVar(&o.Archive.UseLMStudio, "lmstudio", false, "Use model analysis in addition to HTML checks")
		fs.StringVar(&o.Archive.ModelURL, "lmstudio-url", o.Archive.ModelURL, "OpenAI-compatible API base URL")
		fs.StringVar(&o.Archive.Model, "model", o.Archive.Model, "Verification model")
	}
	fs.DurationVar(&o.Archive.RequestTimeout, "request-timeout", o.Archive.RequestTimeout, "HTTP timeout")
	fs.DurationVar(&o.Archive.LockTimeout, "lock-timeout", o.Archive.LockTimeout, "Report lock timeout")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(a.stderr, "Unexpected positional arguments")
		return 1
	}
	o.Archive.DryRun = *dryRun || (redirects && !*update)
	var stats linkaudit.RepairSummary
	var err error
	if redirects {
		stats, err = linkaudit.UpdateRedirects(ctx, o, a.stdout)
	} else {
		stats, err = linkaudit.ArchiveBroken(ctx, o, a.stdout)
	}
	mode := "updated"
	if o.Archive.DryRun {
		mode = "would update"
	}
	fmt.Fprintf(a.stdout, "Links: %d %s, %d failed, %d filtered, %d already processed\n", stats.Updated, mode, stats.Failed, stats.Filtered, stats.Skipped)
	if err != nil {
		return a.auditError(ctx, err)
	}
	if stats.Failed > 0 {
		return 1
	}
	return 0
}

func (a *App) runAuditContent(ctx context.Context, args []string) int {
	defaults, input := auditPaths()
	defaults.Model = "qwen/qwen3-30b-a3b-128k"
	defaults.ModelTimeout = 600 * time.Second
	defaults.MaxContent = 8000
	o := linkaudit.ContentOptions{ContentDir: defaults.ContentDir, Report: filepath.Join(filepath.Dir(filepath.Dir(input)), "hugo_post_audit_report.json"), Resume: true, UseModel: true, Delay: 2 * time.Second, Model: defaults}
	fs := flag.NewFlagSet("audit-content", flag.ContinueOnError)
	fs.SetOutput(a.stdout)
	fs.StringVar(&o.ContentDir, "content-dir", o.ContentDir, "Hugo post bundle directory")
	fs.StringVar(&o.Report, "report", o.Report, "Resumable content audit JSON array")
	fs.StringVar(&o.Model.ModelURL, "lmstudio-url", o.Model.ModelURL, "OpenAI-compatible API base URL")
	fs.StringVar(&o.Model.Model, "model", o.Model.Model, "Content audit model")
	fs.DurationVar(&o.Model.ModelTimeout, "model-timeout", o.Model.ModelTimeout, "Timeout per model request")
	fs.IntVar(&o.Model.Retries, "model-retries", o.Model.Retries, "Model request attempts")
	fs.IntVar(&o.Model.MaxContent, "max-content-length", o.Model.MaxContent, "Maximum characters sent per post")
	fs.DurationVar(&o.Delay, "request-delay", o.Delay, "Delay between posts")
	fs.BoolVar(&o.Resume, "resume", true, "Skip already audited posts")
	posts := fs.String("posts", "", "Comma-separated post slugs")
	noModel := fs.Bool("no-lmstudio", false, "Use local conservative WordPress shortcode detection")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if fs.NArg() > 0 {
		fmt.Fprintln(a.stderr, "Unexpected positional arguments")
		return 1
	}
	o.Posts = parseSlugs(*posts)
	o.UseModel = !*noModel
	results, err := linkaudit.AuditContent(ctx, o, a.stdout)
	if err != nil {
		return a.auditError(ctx, err)
	}
	issues, failures := 0, 0
	for _, r := range results {
		if r.HasIssue {
			issues++
		}
		if r.Error != "" {
			failures++
		}
	}
	fmt.Fprintf(a.stdout, "Content audit: %d/%d posts with potential issues; %d model errors\n", issues, len(results), failures)
	if failures > 0 {
		return 1
	}
	return 0
}
