## Interactive terminal interface

Build the executable as above, then launch from the directory containing your
`config.yaml` (relative paths and environment variables work as in the CLI):

```bash
./bin/pipeline tui
```

The gocui interface exposes all pipeline commands with a command menu, editable
arguments, automatic help for the selected command, and live combined output. It requires a terminal of at least 72 × 18.

- **↑/↓** or **j/k**: select a command in the menu.
- **h**: refresh the selected command's help (also shown automatically).
- **Tab**: switch between the menu and arguments; **Enter** opens arguments.
- **Ctrl-R**: run the selected command with the displayed arguments.
- **Page Up / Page Down**: scroll output; return to the bottom to follow activity.
- **Ctrl-X**: cancel the running command, allowing five seconds for cleanup.
- **Ctrl-C**: cancel any running command and exit.

Arguments are remembered per command for the session. Quote paths or values with
spaces, for example `--posts my-post --config "my config.yaml"`. Arguments are
passed directly to the CLI without shell expansion. One command runs at a time;
output follows the latest activity and retains the most recent 128 KiB. Commands
have their usual effects, so use their `--dry-run` or staging flags when desired.
