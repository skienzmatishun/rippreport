package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rippreport/go-content-pipeline/internal/archive"
)

// Resolve defaults from either the site root or the Go module directory.
func archiveDefaults() archive.Options {
	o := archive.DefaultOptions()
	if _, err := os.Stat("../content/p"); err == nil {
		o.Report = filepath.Join("..", o.Report)
		o.Progress = filepath.Join("..", o.Progress)
		o.ContentDir = filepath.Join("..", o.ContentDir)
		o.BackupDir = filepath.Join("..", o.BackupDir)
		o.ScreenshotDir = filepath.Join("..", o.ScreenshotDir)
	}
	if _, err := os.Stat(o.Report); os.IsNotExist(err) {
		candidate := filepath.Join(filepath.Dir(o.Report), "broken_links_report.json")
		if _, err := os.Stat(candidate); err == nil {
			o.Report = candidate
		}
	}
	return o
}

func (a *App) runMissingPDFs(args []string) int {
	o := archiveDefaults()
	fs := flag.NewFlagSet("missing-pdfs", flag.ContinueOnError)
	fs.SetOutput(a.stdout)
	report := fs.String("report", o.Report, "Archive-fix or broken-links audit JSON")
	output := fs.String("output", filepath.Join(filepath.Dir(o.Report), "missing_pdfs.json"), "Output JSON path")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(a.stderr, "Unexpected positional arguments")
		return 1
	}
	if sameFilePath(*report, *output) {
		fmt.Fprintln(a.stderr, "Output must differ from the source report")
		return 1
	}
	count, err := archive.ExportMissingPDFs(*report, *output)
	if err != nil {
		fmt.Fprintf(a.stderr, "PDF export failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(a.stdout, "Exported %d missing PDFs to %s\n", count, *output)
	return 0
}

func sameFilePath(a, b string) bool { a, _ = filepath.Abs(a); b, _ = filepath.Abs(b); return a == b }

func (a *App) runVerifyArchives(ctx context.Context, args []string) int {
	o := archiveDefaults()
	fs := flag.NewFlagSet("verify-archives", flag.ContinueOnError)
	fs.SetOutput(a.stdout)
	fs.StringVar(&o.Report, "report", o.Report, "Archive-fix or broken-links audit JSON (updated in place)")
	fs.StringVar(&o.Progress, "progress", o.Progress, "Resumable progress JSON")
	fs.StringVar(&o.ContentDir, "content-dir", o.ContentDir, "Hugo post bundle directory")
	fs.StringVar(&o.BackupDir, "backup-dir", o.BackupDir, "Verified original content backups")
	fs.StringVar(&o.ScreenshotDir, "screenshot-dir", o.ScreenshotDir, "Browser screenshot directory")
	fs.StringVar(&o.Python, "python", o.Python, "Python interpreter with Playwright installed")
	fs.StringVar(&o.ModelURL, "lmstudio-url", o.ModelURL, "OpenAI-compatible API base URL")
	fs.StringVar(&o.Model, "model", o.Model, "Text/vision verification model")
	fs.StringVar(&o.CDXURL, "cdx-url", o.CDXURL, "Wayback CDX endpoint")
	fs.StringVar(&o.SearchURL, "search-url", o.SearchURL, "Google search endpoint")
	fs.BoolVar(&o.DryRun, "dry-run", false, "Preview without writing content, backups, reports or progress")
	noBrowser := fs.Bool("no-playwright", false, "Use HTTP only; no Python dependency")
	noModel := fs.Bool("no-lmstudio", false, "Use HTML/content-type verification only")
	noGoogle := fs.Bool("no-google", false, "Disable network Google fallback searches")
	fs.BoolVar(&o.Vision, "vision", o.Vision, "Analyze screenshots when a vision model is available")
	fs.BoolVar(&o.VerifyGoogle, "verify-google-results", o.VerifyGoogle, "Verify relevance of canonical post search results")
	fs.BoolVar(&o.FixCDN, "fix-malformed-cdn", o.FixCDN, "Repair malformed CDN destinations")
	fs.BoolVar(&o.ConvertTags, "convert-tags", o.ConvertTags, "Convert tag links to site searches")
	fs.BoolVar(&o.FixMissingPath, "fix-missing-p-path", o.FixMissingPath, "Repair missing /p/ paths")
	fs.BoolVar(&o.VerifyPaywall, "verify-paywall", o.VerifyPaywall, "Check accessibility and article relevance")
	fs.BoolVar(&o.CleanupScreenshots, "cleanup-screenshots", o.CleanupScreenshots, "Delete screenshots after analysis")
	fs.DurationVar(&o.RequestTimeout, "request-timeout", o.RequestTimeout, "HTTP timeout")
	fs.DurationVar(&o.BrowserTimeout, "browser-timeout", o.BrowserTimeout, "Browser navigation timeout")
	fs.DurationVar(&o.ModelTimeout, "lmstudio-timeout", o.ModelTimeout, "Timeout per model attempt")
	fs.DurationVar(&o.Wait, "browser-wait", o.Wait, "Wait for dynamic content after navigation")
	fs.DurationVar(&o.SearchDelay, "search-delay", o.SearchDelay, "Minimum interval between searches")
	fs.DurationVar(&o.SearchJitter, "search-jitter", o.SearchJitter, "Additional random search delay")
	fs.DurationVar(&o.LockTimeout, "lock-timeout", o.LockTimeout, "Timeout acquiring report/progress locks")
	fs.IntVar(&o.MaxSnapshots, "max-snapshots", o.MaxSnapshots, "Maximum previous snapshots to try")
	fs.IntVar(&o.MaxResults, "max-results", o.MaxResults, "Maximum Google results")
	fs.IntVar(&o.MaxContent, "max-content-length", o.MaxContent, "Maximum characters sent to model")
	fs.IntVar(&o.Retries, "lmstudio-retries", o.Retries, "Model attempts with exponential backoff")
	logFile := fs.String("log-file", "", "Optional log file (also writes to console)")
	logLevel := fs.String("log-level", "INFO", "Log level: DEBUG, INFO, WARNING, ERROR")
	missingPDFs := fs.String("missing-pdfs-output", "", "Export remaining unresolved PDFs after verification")
	restore := fs.String("restore-backup", "", "Restore a backup manifest instead of verification")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if fs.NArg() != 0 {
		fmt.Fprintln(a.stderr, "Unexpected positional arguments")
		return 1
	}
	switch *logLevel {
	case "DEBUG", "INFO", "WARNING", "ERROR":
	default:
		fmt.Fprintln(a.stderr, "Invalid log level")
		return 1
	}
	o.UsePlaywright = !*noBrowser
	o.UseLMStudio = !*noModel
	o.Google = !*noGoogle
	if err := o.Validate(); err != nil {
		fmt.Fprintf(a.stderr, "Invalid archive options: %v\n", err)
		return 1
	}
	if sameFilePath(o.Report, o.Progress) || (*missingPDFs != "" && (sameFilePath(o.Report, *missingPDFs) || sameFilePath(o.Progress, *missingPDFs))) {
		fmt.Fprintln(a.stderr, "Report, progress and PDF output must use distinct files")
		return 1
	}
	if *restore != "" {
		if o.DryRun {
			fmt.Fprintln(a.stderr, "--restore-backup cannot be combined with --dry-run")
			return 1
		}
		if err := archive.RestoreBackup(*restore); err != nil {
			fmt.Fprintf(a.stderr, "Restore failed: %v\n", err)
			return 1
		}
		fmt.Fprintln(a.stdout, "Backup restored")
		return 0
	}
	var log io.Writer = a.stdout
	if *logLevel == "WARNING" || *logLevel == "ERROR" {
		log = io.Discard
	}
	if *logFile != "" {
		if sameFilePath(*logFile, o.Report) || sameFilePath(*logFile, o.Progress) {
			fmt.Fprintln(a.stderr, "Log file must differ from report and progress")
			return 1
		}
		f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(a.stderr, "Cannot open log: %v\n", err)
			return 1
		}
		defer f.Close()
		log = io.MultiWriter(log, f)
	}
	v := archive.New(o, log)
	summary, err := v.Run(ctx)
	fmt.Fprintf(a.stdout, "Archives: %d replaced, %d verified, %d would replace, %d failed, %d skipped (%d total)\n", summary.Replaced, summary.Verified, summary.WouldReplace, summary.Failed, summary.Skipped, summary.Total)
	if err != nil {
		fmt.Fprintf(a.stderr, "Archive verification failed: %v\n", err)
		if ctx.Err() != nil {
			return 130
		}
		return 1
	}
	if *missingPDFs != "" {
		if o.DryRun {
			fmt.Fprintln(a.stdout, "Dry run: PDF export skipped; use missing-pdfs to export the current report")
		} else {
			count, err := archive.ExportMissingPDFs(o.Report, *missingPDFs)
			if err != nil {
				fmt.Fprintf(a.stderr, "PDF export failed: %v\n", err)
				return 1
			}
			fmt.Fprintf(a.stdout, "Exported %d missing PDFs to %s\n", count, *missingPDFs)
		}
	}
	if summary.Failed > 0 {
		return 1
	}
	return 0
}
