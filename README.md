<p align="center">
<strong>セッション</strong>
</p>

<p align="center">
  A Go CLI tool for managing predefined TMUX sessions with workspaces, windows, and panels.
</p>

## Features

- **Workspace Organization**: Group related sessions into workspaces (work, personal, labs)
- **Custom Windows & Panels**: Define any number of windows with custom layouts and commands
- **Working Directory Support**: Save and restore working directories per panel (supports `~`, `$HOME`, `${HOME}`)
- **Interactive Session Selection**: Choose which sessions to spawn from a workspace
- **Safe Session Management**: Automatic fallback when killing sessions
- **Save Current Sessions**: Capture your current TMUX setup with one command
- **Diagnostics**: `sesh doctor` validates config and detects tmux/session mismatches
- **Configuration in YAML**: Easy to read and edit configuration format

## Installation

### From Source

```bash
cd sesh
make install
```

Or manually:

```bash
cd sesh
go build -o sesh ./cmd/sesh
sudo mv sesh /usr/local/bin/
```

### From Release

Download a pre-built binary from the [releases page](https://github.com/yourusername/sesh/releases):

```bash
# macOS (Apple Silicon)
curl -L -o sesh https://github.com/yourusername/sesh/releases/latest/download/sesh-darwin-arm64
chmod +x sesh
sudo mv sesh /usr/local/bin/

# Linux (amd64)
curl -L -o sesh https://github.com/yourusername/sesh/releases/latest/download/sesh-linux-amd64
chmod +x sesh
sudo mv sesh /usr/local/bin/
```

Verify the checksum against `checksums.txt` in the release.

## Quick Start

```bash
# Create a workspace
sesh workspace create work --description "Work projects"

# Create a session in a workspace
sesh create myproject -w "editor" -w "terminal:2" --workspace work

# Spawn all sessions in a workspace (interactive selection)
sesh workspace spawn work

# View all workspaces and sessions
sesh list

# Edit configuration manually
sesh edit
```

## Configuration

Sesh stores configuration in `~/.config/sesh/sessions.yaml`:

```yaml
workspaces:
  - name: work
    description: Work-related projects
    sessions:
      - name: backend-api
        windows:
          - name: editor
            workdir: ~/work/backend-api
            layout: main-vertical
            panels:
              - command: nvim .
                workdir: ~/work/backend-api
              - command: git status
                workdir: ~/work/backend-api
          - name: terminal
            panels:
              - {}

  - name: personal
    description: Personal projects
    sessions:
      - name: dotfiles
        windows:
          - name: config
            workdir: ~/.config
            panels:
              - command: nvim
```

## Commands

### Session Management

| Command                                                   | Description                                                              |
| --------------------------------------------------------- | ------------------------------------------------------------------------ |
| `sesh create [name] -w [windows] --workspace [workspace]` | Create a session in a workspace                                          |
| `sesh spawn [name]`                                       | Spawn a specific session                                                 |
| `sesh kill [name]`                                        | Kill a running session (works for untracked sessions too)                |
| `sesh kill --all`                                         | Kill all running sessions                                                |
| `sesh kill --workspace [name]`                            | Kill all sessions in a workspace                                         |
| `sesh kill --dry-run`                                     | Preview what `kill` would terminate                                      |
| `sesh delete [name]`                                      | Delete a session definition                                              |
| `sesh move [name] --workspace [workspace]`                | Move a session to a workspace                                            |
| `sesh move [name] --standalone`                           | Move a session to standalone (remove from workspace)                     |
| `sesh rename [old] [new]`                                 | Rename a session (updates config + tmux)                                 |
| `sesh show [name]`                                        | Show session details (windows, panels, status)                           |
| `sesh clone [name] [new-name]`                            | Clone a session definition                                               |
| `sesh add-window [session] [definition]`                  | Add a window to a session                                                |
| `sesh remove-window [session] [window]`                   | Remove a window from a session                                           |
| `sesh add-panel [session] [window]`                       | Add a panel to a window                                                  |
| `sesh remove-panel [session] [window] [index]`            | Remove a panel by index                                                  |
| `sesh list`                                               | List all sessions grouped by workspace (shows untracked sessions with ◌) |
| `sesh list --running`                                     | Show only running sessions                                               |
| `sesh list --workspace [name]`                            | Show only sessions in a workspace                                        |
| `sesh list --standalone`                                  | Show only standalone sessions                                            |

### Workspace Management

| Command                                             | Description                                      |
| --------------------------------------------------- | ------------------------------------------------ |
| `sesh workspace create [name] --description "desc"` | Create a new workspace                           |
| `sesh workspace list`                               | List all workspaces                              |
| `sesh workspace show [workspace]`                   | Show workspace details and session status        |
| `sesh workspace spawn [workspace]`                  | Interactive session selection                    |
| `sesh workspace spawn [workspace] --all`            | Spawn all sessions without prompting             |
| `sesh workspace kill [workspace]`                   | Kill all sessions in workspace                   |
| `sesh workspace delete [workspace]`                 | Delete workspace (asks to kill running sessions) |
| `sesh workspace rename [old] [new]`                 | Rename a workspace                               |

### Diagnostics

| Command                                         | Description                                                                          |
| ----------------------------------------------- | ------------------------------------------------------------------------------------ |
| `sesh doctor`                                   | Check config validation and tmux inconsistencies (colored output, exits 1 on issues) |
| `sesh config validate`                          | Validate config file for errors and warnings                                         |
| `sesh completion [bash\|zsh\|fish\|powershell]` | Generate shell completion script                                                     |

### Configuration

| Command                             | Description                                                     |
| ----------------------------------- | --------------------------------------------------------------- |
| `sesh save --workspace [workspace]` | Save a new untracked tmux session into config (initial capture) |
| `sesh save --dry-run`               | Preview what `save` would do without writing                    |
| `sesh sync [session]`               | Update an existing config session to match current tmux state   |
| `sesh sync --all`                   | Update all running config sessions to match current tmux state  |
| `sesh sync --dry-run`               | Preview what `sync` would update without writing                |
| `sesh delete [name] --dry-run`      | Preview what `delete` would remove                              |
| `sesh edit`                         | Open config in default editor                                   |
| `sesh config`                       | Show config file path                                           |

**`save` vs `sync`:**

- **`save`** — Use when you manually created a tmux session and want to add it to sesh config for the first time
- **`sync`** — Use when a session already exists in config and you've made changes in tmux (new windows, panels, etc.) that you want reflected in the config
- **Hooks** — Add tmux hooks to auto-run `sync` on window/pane changes (see Auto-Sync section below)

## Window Definition

When creating sessions, use the `-w` flag:

```bash
# Format: name or name:panel_count
sesh create myproject \
  -w "editor:2" \
  -w "terminal" \
  -w "logs:3" \
  --workspace work
```

## Layouts

Available TMUX layouts:

- `even-horizontal` - Panels side by side
- `even-vertical` - Panels stacked
- `main-horizontal` - Main panel on top
- `main-vertical` - Main panel on left
- `tiled` - Grid arrangement

## Workflow Examples

### Setup Work Environment

```bash
# Create workspace and sessions
sesh workspace create work --description "Work projects"
sesh create backend -w "editor:2" -w "server" --workspace work
sesh create frontend -w "dev" -w "build" --workspace work

# Start entire work environment (interactive selection)
sesh workspace spawn work
# ? Select sessions to spawn from workspace 'work':
#   [✓] backend (2 windows)
#   [✓] frontend (2 windows)

# When done
sesh workspace kill work
# or
sesh workspace delete work
```

### Save Current Session

```bash
# From inside a TMUX session
cd ~/myproject
sesh save --workspace work
# Captures current structure, windows, and working directories
```

### Move Sessions

```bash
# Move a standalone session into a workspace
sesh move myproject --workspace work

# Move a session between workspaces
sesh move backend --workspace staging

# Move a session out of a workspace to standalone
sesh move backend --standalone
```

### Rename Sessions

```bash
# Rename a session (updates both config and running tmux session)
sesh rename old-name new-name
```

### Show Session Details

```bash
# Show detailed info about a session
sesh show backend
# Session: backend (workspace: work)
# Status:  running
# Windows: 2
#
#   Window: editor
#     WorkDir: ~/work/backend
#     Layout:  main-vertical
#     Panel 1:
#       Command: nvim .
#       WorkDir: ~/work/backend
#     Panel 2: (default shell)
```

### Diagnose Issues

```bash
# Check for config/tmux mismatches
sesh doctor
# Output:
#   Config sessions not running:
#     - work/backend
#
#   Running TMUX sessions not in config:
#     - work/backend-api
#
#   Possible renames (use 'sesh rename' to fix):
#     - work/backend -> work/backend-api
```

The `list` command also shows untracked sessions (running in tmux but not in config) with the ◌ indicator:

```
[standalone]
  ● myproject (2 windows)

[workspaces]
  work: Work projects
    ○ backend (2 windows)
    ● frontend (1 windows)

[untracked]
  ◌ scratch
```

### Manage Workspaces

```bash
# View all workspaces
sesh workspace list

# Create empty workspace
sesh workspace create experiments

# Delete a workspace (confirms before killing running sessions)
sesh workspace delete old-project
```

## Environment Variable Expansion

Working directory paths in config support `~`, `$HOME`, `${HOME}`, and other environment variables:

```yaml
windows:
  - name: editor
    workdir: ~/work/project # expands to /home/user/work/project
    panels:
      - workdir: $HOME/projects # expands using $HOME
```

## Shell Completion

```bash
# Bash
eval "$(sesh completion bash)"

# Zsh
eval "$(sesh completion zsh)"

# Fish
sesh completion fish | source
```

## Auto-Sync with Tmux Hooks

Add these hooks to your `~/.tmux.conf` to automatically sync session structure back to config when windows or panels change:

```bash
# Auto-sync when creating windows/panels
set-hook -g after-new-window 'run-shell "sesh sync #S > /dev/null 2>&1 || true"'
set-hook -g after-split-window 'run-shell "sesh sync #S > /dev/null 2>&1 || true"'

# Auto-sync when closing windows/panels
set-hook -g after-kill-pane 'run-shell "sesh sync #S > /dev/null 2>&1 || true"'
set-hook -g pane-exited 'run-shell "sesh sync #S > /dev/null 2>&1 || true"'
```

**How it works:**

- `#S` expands to the current session name
- `sesh sync` introspects the running tmux state and updates the config file
- `> /dev/null 2>&1 || true` suppresses output and prevents errors from blocking tmux
- Only sessions already defined in config are synced — untracked sessions are ignored

**Note:** This syncs on every structural change. If you make rapid changes (e.g., creating and destroying many windows quickly), consider running `sesh sync --all` manually after your session settles instead.

## Verbose Mode

Pass `-v` or `--verbose` to log tmux commands being executed:

```bash
sesh -v spawn myproject
# 2025/04/17 10:00:00 [tmux] tmux has-session -t =myproject
# 2025/04/17 10:00:00 [tmux] tmux new-session -d -s myproject -n editor
```

## Safety Features

- **Automatic Fallback**: When killing a session you're attached to, automatically switches to another available session
- **Confirmation Prompts**: `workspace delete` asks to kill running sessions
- **Error Propagation**: All tmux operations propagate errors; non-fatal failures (panel commands, directory changes, layout changes) are reported as warnings instead of silently ignored
- **Session Introspection Warnings**: When saving a session, per-panel failures (command detection, working directory) are collected and reported rather than discarded
- **Session Existence Check**: Prevents duplicate sessions

## Testing

The project includes a comprehensive test suite with full mocking support for tmux commands, allowing tests to run without a tmux server.

```bash
make test       # Run all tests
make lint       # Run fmt and vet
```

### Test Coverage

| Package             | Tests | Description                                                                               |
| ------------------- | ----- | ----------------------------------------------------------------------------------------- |
| `internal/config`   | 62    | Config CRUD, session/workspace management, move/rename/clone/window/panel/sync operations |
| `internal/tmux`     | 69    | Session lifecycle, introspection, query, client, command runner                           |
| `internal/doctor`   | 12    | Config/tmux mismatch detection, rename suggestions                                        |
| `internal/expand`   | 7     | Path expansion (`~`, `$HOME`, env vars)                                                   |
| `internal/output`   | 6     | Output helpers (Info/Warn/Error)                                                          |
| `internal/validate` | 11    | Config validation (duplicates, empty names, invalid layouts)                              |
| `internal/editor`   | 3     | Editor detection and configuration                                                        |
| `internal/parser`   | 6     | Window definition parsing                                                                 |
| `pkg/models`        | 6     | Data structures and YAML marshaling                                                       |

### Architecture

The tmux client uses a `CommandRunner` interface for dependency injection, making all tmux operations testable without a running tmux server:

```go
type CommandRunner interface {
    Run(name string, args ...string) error
    Output(name string, args ...string) (string, error)
}
```

`NewClient()` uses a `RealRunner` (exec.Command) for production, while `NewClientWithRunner()` accepts any implementation for testing.

## Makefile Targets

```bash
make build      # Build the binary
make install    # Build and install to /usr/local/bin
make uninstall  # Remove from /usr/local/bin
make clean      # Remove build artifacts
make test       # Run tests
make lint       # Run fmt and vet
make tidy       # Run go mod tidy
make release    # Build cross-platform binaries (macOS, Linux)
make help       # Show all targets
```

## Releasing

To create a new release:

1. Update `version` in `cmd/sesh/root.go`
2. Commit and tag:
   ```bash
   git add .
   git commit -m "release: v0.2.0"
   git tag v0.2.0
   git push origin v0.2.0
   ```
3. GitHub Actions automatically builds binaries and creates a release
4. Or build locally: `make release` (outputs to `build/release/`)

## License

MIT
