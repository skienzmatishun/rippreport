// Package tui provides an interactive front end for the pipeline CLI.
package tui

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
)

var commands = []string{
	"generate-related", "generate-hooks", "generate-frontpage", "pipeline",
	"refine-hooks", "dedup-hooks", "apply-staged", "clear-staged",
	"fetch-transcript", "summarize-transcript", "generate-alt", "social-post",
	"new-post", "email-post", "repost", "send-newsletter",
	"check-links", "analyze-broken-links", "archive-links", "update-redirected-links",
	"audit-content", "verify-archives", "missing-pdfs",
	"scan-hook-articles", "scan-hook-authors", "scan-hook-phrases",
}

type app struct {
	help        string
	helpCommand string
	helpCache   map[string]string
	scroll      int
	selected    int
	args        map[int]string
	output      string
	status      string
	running     bool
	cancel      context.CancelFunc
	ctx         context.Context
	executable  string
	workers     sync.WaitGroup
}

// Run opens the interface. Commands execute in the current working directory,
// inheriting the environment and the existing CLI's configuration defaults.
func Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("tui", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("usage: pipeline tui")
	}
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return fmt.Errorf("open terminal interface: %w", err)
	}
	defer g.Close()
	a := &app{ctx: ctx, executable: executable, args: make(map[int]string), status: "Ready", helpCache: make(map[string]string), output: "Select a command to see its options, then Tab to edit arguments.\nCtrl-R runs the displayed command.\n"}
	defer func() {
		if a.cancel != nil {
			a.cancel()
		}
		a.workers.Wait()
	}()
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan
	g.SetManagerFunc(a.layout)
	bindings := []struct {
		view string
		key  interface{}
		fn   func(*gocui.Gui, *gocui.View) error
	}{
		{"", gocui.KeyCtrlC, func(*gocui.Gui, *gocui.View) error { return gocui.ErrQuit }},
		{"", gocui.KeyCtrlR, func(g *gocui.Gui, v *gocui.View) error { return a.start(g, false) }},
		{"", gocui.KeyCtrlX, func(g *gocui.Gui, v *gocui.View) error {
			if a.cancel != nil {
				a.cancel()
				a.status = "Cancelling…"
			}
			return nil
		}},
		{"", gocui.KeyTab, func(g *gocui.Gui, v *gocui.View) error {
			name := "arguments"
			if v != nil && v.Name() == name {
				name = "commands"
			}
			_, err := g.SetCurrentView(name)
			g.Cursor = name == "arguments"
			return err
		}},
		{"", gocui.KeyPgup, a.scrollOutput(1)},
		{"", gocui.KeyPgdn, a.scrollOutput(-1)},
		{"commands", gocui.KeyArrowDown, a.move(1)}, {"commands", gocui.KeyArrowUp, a.move(-1)},
		{"commands", 'j', a.move(1)}, {"commands", 'k', a.move(-1)},
		{"commands", 'h', func(g *gocui.Gui, v *gocui.View) error { a.showHelp(g, true); return nil }},
		{"commands", gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
			_, err := g.SetCurrentView("arguments")
			g.Cursor = true
			return err
		}},
	}
	for _, b := range bindings {
		if err := g.SetKeybinding(b.view, b.key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (a *app) layout(g *gocui.Gui) error {
	w, h := g.Size()
	if w < 72 || h < 18 {
		for _, name := range []string{"commands", "arguments", "help", "output", "status"} {
			_ = g.DeleteView(name)
		}
		v, err := g.SetView("small", 0, 0, w-1, h-1)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		fmt.Fprint(v, "Resize terminal to at least 72 × 18. Ctrl-C exits.")
		return nil
	}
	_ = g.DeleteView("small")
	left := 28
	menu, err := g.SetView("commands", 0, 0, left, h-4)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	menu.Title = " Pipeline commands "
	menu.Highlight = true
	menu.SelFgColor = gocui.ColorBlack
	menu.SelBgColor = gocui.ColorCyan
	menu.Clear()
	for _, cmd := range commands {
		fmt.Fprintln(menu, cmd)
	}
	_, height := menu.Size()
	origin := 0
	if a.selected >= height {
		origin = a.selected - height + 1
	}
	_ = menu.SetOrigin(0, origin)
	_ = menu.SetCursor(0, a.selected-origin)
	input, err := g.SetView("arguments", left+1, 0, w-1, 3)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	input.Title = " Arguments · " + commands[a.selected] + " "
	input.Editable = true
	if err == gocui.ErrUnknownView {
		fmt.Fprint(input, a.args[a.selected])
	}
	if a.helpCommand == "" {
		a.showHelp(g, false)
	}
	help, err := g.SetView("help", left+1, 4, w-1, h/2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	help.Title = " Help · " + commands[a.selected] + " "
	help.Wrap = true
	help.Clear()
	fmt.Fprint(help, a.help)
	out, err := g.SetView("output", left+1, h/2+1, w-1, h-4)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	out.Title = " Output "
	out.Wrap = true
	out.Autoscroll = a.scroll == 0
	out.Clear()
	fmt.Fprint(out, a.output)
	if a.scroll > 0 {
		_, height := out.Size()
		origin := len(out.ViewBufferLines()) - height - a.scroll
		if origin < 0 {
			origin = 0
		}
		_ = out.SetOrigin(0, origin)
	}
	status, err := g.SetView("status", 0, h-3, w-1, h-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	status.Frame = false
	status.Clear()
	fmt.Fprintf(status, "%s | Tab focus · ↑/↓ select · h help · Ctrl-R run · Ctrl-X cancel · Ctrl-C quit", a.status)
	if g.CurrentView() == nil {
		_, err = g.SetCurrentView("commands")
		return err
	}
	return nil
}

func (a *app) move(delta int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if input, err := g.View("arguments"); err == nil {
			a.args[a.selected] = strings.TrimSpace(input.Buffer())
		}
		a.selected = (a.selected + delta + len(commands)) % len(commands)
		if input, err := g.View("arguments"); err == nil {
			input.Clear()
			_ = input.SetOrigin(0, 0)
			_ = input.SetCursor(0, 0)
			fmt.Fprint(input, a.args[a.selected])
		}
		a.showHelp(g, false)
		return nil
	}
}

// guiWriter serializes output from both child streams onto the UI event loop.
type guiWriter struct {
	g *gocui.Gui
	a *app
}

func (w guiWriter) Write(p []byte) (int, error) {
	s := string(p)
	w.g.Update(func(*gocui.Gui) error {
		w.a.output += s
		// Bound retained output so long batch jobs do not consume unlimited memory.
		if len(w.a.output) > 128*1024 {
			w.a.output = w.a.output[len(w.a.output)-128*1024:]
		}
		return nil
	})
	return len(p), nil
}

func (a *app) start(g *gocui.Gui, help bool) error {
	if a.running {
		a.status = "A command is running; Ctrl-X cancels it"
		return nil
	}
	input, err := g.View("arguments")
	if err != nil {
		return nil
	}
	raw := strings.TrimSpace(input.Buffer())
	a.args[a.selected] = raw
	args, err := splitArgs(raw)
	if err != nil {
		a.status = err.Error()
		return nil
	}
	if help {
		args = []string{"--help"}
	}
	args = append([]string{commands[a.selected]}, args...)
	ctx, cancel := context.WithCancel(a.ctx)
	a.cancel = cancel
	a.running = true
	a.scroll = 0
	a.status = "Running " + commands[a.selected]
	a.output = "$ pipeline " + strings.Join(args, " ") + "\n\n"
	a.workers.Add(1)
	go func() {
		defer a.workers.Done()
		defer cancel()
		cmd := exec.CommandContext(ctx, a.executable, args...)
		// Let the CLI perform its normal graceful cleanup before forcing termination.
		cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
		cmd.WaitDelay = 5 * time.Second
		writer := guiWriter{g, a}
		cmd.Stdout = writer
		cmd.Stderr = writer
		cmd.Stdin = nil
		err := cmd.Run()
		cancelled := ctx.Err() != nil
		g.Update(func(*gocui.Gui) error {
			a.running = false
			a.cancel = nil
			a.status = "Completed"
			if err != nil {
				a.status = "Failed: " + err.Error()
			}
			if cancelled {
				a.status = "Cancelled"
			}
			return nil
		})
	}()
	return nil
}

// splitArgs handles shell-style quoting without invoking a shell or expanding
// variables, substitutions, globs, or redirections.
func splitArgs(s string) ([]string, error) {
	var args []string
	var word strings.Builder
	var quote rune
	escaped, started := false, false
	for _, r := range s {
		if escaped {
			word.WriteRune(r)
			escaped = false
			started = true
			continue
		}
		if r == '\\' && quote != '\'' {
			escaped = true
			started = true
			continue
		}
		if quote != 0 {
			if r == quote {
				quote = 0
			} else {
				word.WriteRune(r)
			}
			continue
		}
		switch r {
		case '\'', '"':
			quote = r
			started = true
		case ' ', '\t', '\n':
			if started {
				args = append(args, word.String())
				word.Reset()
				started = false
			}
		default:
			word.WriteRune(r)
			started = true
		}
	}
	if escaped || quote != 0 {
		return nil, fmt.Errorf("arguments contain an unfinished quote or escape")
	}
	if started {
		args = append(args, word.String())
	}
	return args, nil
}

var _ io.Writer = guiWriter{}

func (a *app) scrollOutput(direction int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, _ *gocui.View) error {
		if v, err := g.View("output"); err == nil {
			_, h := v.Size()
			a.scroll += direction * (h - 1)
			limit := len(v.ViewBufferLines()) - h
			if limit < 0 {
				limit = 0
			}
			if a.scroll > limit {
				a.scroll = limit
			}
			if a.scroll < 0 {
				a.scroll = 0
			}
		}
		return nil
	}
}

// showHelp runs independently of jobs so browsing options never replaces job
// output. Discard stale responses when the selection changes during a request.
func (a *app) showHelp(g *gocui.Gui, refresh bool) {
	command := commands[a.selected]
	a.helpCommand = command
	if help, ok := a.helpCache[command]; ok && !refresh {
		a.help = help
		return
	}
	a.help = "Loading options…"
	a.workers.Add(1)
	go func() {
		defer a.workers.Done()
		ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
		defer cancel()
		output, err := exec.CommandContext(ctx, a.executable, command, "--help").CombinedOutput()
		help := string(output)
		if err != nil {
			help += "\nHelp unavailable: " + err.Error()
		}
		g.Update(func(*gocui.Gui) error {
			a.helpCache[command] = help
			if a.helpCommand == command {
				a.help = help
			}
			return nil
		})
	}()
}
