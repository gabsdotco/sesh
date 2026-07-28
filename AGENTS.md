# AGENTS.md

## Commands

```bash
make build        # go build -o build/sesh ./cmd/sesh
make test         # gotestsum --format testname ./... (falls back to go test -v ./...)
make test-cover   # go test -coverprofile=coverage.out ./... + HTML report
make lint         # go fmt ./... && go vet ./...
make install      # build + sudo cp to /usr/local/bin/sesh
```

Run a single test: `go test ./internal/config -run TestRenameSession -v`

## Architecture

- `cmd/sesh/` — Cobra CLI commands. Each file = one command. Workspace subcommands split into separate files (`workspace_create.go`, `workspace_kill.go`, etc.). Dependencies injected via cobra context (`appFromContext(cmd)` returns `*appContext` with `cfgManager`, `tmuxClient`, `edit`).
- `internal/config/` — Config CRUD. `Manager` reads/writes `~/.config/sesh/sessions.yaml`. Session names are globally unique across workspaces and orphans.
- `internal/tmux/` — Tmux client. `CommandRunner` interface (`Run`, `Output`) enables mocking. `RealRunner` calls `os/exec`; `MockRunner` uses pattern matching for tests. Methods organized by domain: `session.go`, `attach.go`, `workspace.go`, `query.go`, `introspect.go`.
- `internal/doctor/` — Pure logic comparing config vs running sessions; no tmux calls.
- `internal/output/` — `Info()`, `Warn()`, `Error()` helpers. All user output uses these instead of raw `fmt.Printf`. `Warn` and `Error` go to stderr.
- `internal/parser/` — Window definition shorthand parser (`"name:count"`).
- `internal/editor/` — Opens `$EDITOR` (defaults to vim).
- `internal/expand/` — Path expansion: `~`, `$HOME`, `${HOME}`, and `$VAR` in WorkDir fields.
- `internal/validate/` — Config validation: duplicate names, empty names, invalid layouts, missing windows. `FilterBySeverity()` helper for consumers.
- `pkg/models/` — YAML-tagged structs: `Config`, `Workspace`, `Session`, `Window`, `Panel`. `Session.TmuxName()` returns the tmux session name (currently just `Name`).

## Testing Conventions

- All tmux tests use `MockRunner`/`MockResponse` from `internal/tmux/mock.go` (non-test file, exported).
- Mock pattern matching: `MockRunner.Run`/`Output` checks `strings.Contains(cmd, resp.Pattern)` where `cmd` = `"tmux " + args joined by spaces`.
- **Exact-match tmux targets**: `SessionExists` uses `-t =name` (the `=` prefix forces exact match; without it, tmux does prefix matching).
- `NewManagerWithPath(path)` creates a `Manager` pointing at a temp file — use this in tests, not `NewManager()`.
- Pure logic (e.g. `doctor.Diagnose`, `parser.ParseWindowDefinition`, `expand.Path`, `validate.Validate`) is extracted to dedicated packages for testability rather than living in `cmd/`.

## Conventions

- No comments in code unless explicitly requested.
- All non-fatal warnings use `output.Warn()` which writes to stderr with `Warning:` prefix. Fatal errors propagate as `error`.
- `Session.TmuxName()` returns just `Name` — session names must be globally unique.
- Module path is `sesh` (local, not a VCS path).
- `--verbose` / `-v` flag sets `tmuxClient.Verbose = true`, which logs all tmux commands via `log.Printf("[tmux] ...")`.
- Command output (success messages, listings) goes to stdout. Diagnostic/warning output goes to stderr via `output.Warn()` / `output.Error()`.

## Gotchas

- Command dependencies are in cobra context, not package-level globals. Use `app := appFromContext(cmd)` at the start of `RunE` functions.
- `MockRunner` patterns match greedily by substring. For methods that call `SessionExists` multiple times with different names, use specific patterns like `"has-session -t =old-name"` and `"has-session -t =new-name"` to avoid false matches.