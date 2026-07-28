<p align="center">
<strong>セッション</strong>
</p>

<p align="center">
  Manage predefined TMUX sessions with workspaces, windows, and panels.
</p>

## Install

From source:

```bash
cd sesh && make install
```

From release:

```bash
curl -L -o sesh https://github.com/gabsdotco/sesh/releases/latest/download/sesh-darwin-arm64
chmod +x sesh && sudo mv sesh /usr/local/bin/
```

## Quick Start

```bash
sesh workspace create work
sesh create myproject -w "editor" -w "terminal:2" --workspace work
sesh workspace spawn work
```

## Commands

| Command | Description |
|---------|-------------|
| `sesh create [name] -w [def]` | Create a session (`-w` accepts `name:count`) |
| `sesh spawn [name]` | Spawn a session |
| `sesh kill [name]` | Kill a session |
| `sesh kill --all` | Kill all running sessions |
| `sesh delete [name]` | Delete from config |
| `sesh rename [old] [new]` | Rename a session |
| `sesh move [name] --workspace [ws]` | Move to workspace |
| `sesh move [name] --standalone` | Move to standalone |
| `sesh clone [name] [new]` | Clone a session |
| `sesh show [name]` | Show session details |
| `sesh list` | List sessions (● running, ○ not running, ◌ untracked) |
| `sesh list --running` | Only running sessions |
| `sesh workspace create [name]` | Create workspace |
| `sesh workspace spawn [ws]` | Spawn workspace sessions |
| `sesh workspace kill [ws]` | Kill workspace sessions |
| `sesh workspace delete [ws]` | Delete workspace |
| `sesh save` | Capture current tmux session into config |
| `sesh save --dry-run` | Preview what save would do |
| `sesh sync [name]` | Update config to match running tmux state |
| `sesh sync --all` | Sync all running sessions |
| `sesh sync --dry-run` | Preview what sync would update |
| `sesh doctor` | Check config/tmux inconsistencies |
| `sesh edit` | Edit config in $EDITOR |
| `sesh config` | Show config path |
| `sesh completion [bash\|zsh\|fish]` | Generate shell completions |

**Save vs Sync:**
- `save` — Add a new untracked tmux session to config for the first time
- `sync` — Update an existing config session after making changes in tmux

## Config

Stored in `~/.config/sesh/sessions.yaml`.

```yaml
workspaces:
  - name: work
    sessions:
      - name: backend
        windows:
          - name: editor
            workdir: ~/work/backend
            panels:
              - command: nvim .
              - {}
sessions:
  - name: scratch
    windows:
      - name: terminal
        panels:
          - {}
```

Working directories support `~`, `$HOME`, and `$VAR` expansion.

## Auto-Sync

Add to `~/.tmux.conf` to keep config in sync automatically:

```bash
set-hook -g after-new-window 'run-shell "sesh sync #S > /dev/null 2>&1 || true"'
set-hook -g after-kill-pane 'run-shell "sesh sync #S > /dev/null 2>&1 || true"'
```

## Development

```bash
make test       # Run tests (183 tests)
make lint       # Format and vet
make release    # Build cross-platform binaries
```

## License

MIT
