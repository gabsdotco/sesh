# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Session save command to capture current tmux session into config
- Session sync command to update existing config sessions from running tmux state
- `sesh sync --all` to sync all running sessions at once
- `--dry-run` flag for `save`, `sync`, `delete`, and `kill` commands
- Doctor command to validate config and detect tmux/session mismatches
- Config validation with `sesh config validate`
- Window and panel CRUD commands (`add-window`, `remove-window`, `add-panel`, `remove-panel`)
- Session clone command
- Workspace rename and show commands
- Batch spawn with `sesh workspace spawn --all`
- Batch kill with `sesh kill --all` and `sesh kill --workspace`
- List filters (`--running`, `--workspace`, `--standalone`)
- Shell completion generation (`sesh completion`)
- Verbose mode (`-v`) to log tmux commands
- Environment variable expansion in working directories (`~`, `$HOME`, `${HOME}`)
- Config backup before write (`.bak` file)
- `sesh --version` flag

### Changed
- `SessionExists` uses exact matching (`-t =name`) to avoid prefix match bugs
- Warning output goes to stderr via `output.Warn()` helper
- `GetFullName()` renamed to `TmuxName()` for clarity
- Workspace subcommands split into separate files
- Dependencies injected via cobra context instead of package globals

### Fixed
- `MockRunner` pattern matching for `SessionExists` with exact names
- `sesh list --standalone` showing workspace sessions
- `sesh list --running` showing empty workspaces
- `make build` failing due to missing Makefile variables

## [0.1.0] - 2025-04-17

### Added
- Initial release with core session management
- Workspace support with interactive session selection
- Create, spawn, kill, delete, move, rename commands
- YAML configuration with `sesh edit`
- Makefile with build, test, lint targets
- Comprehensive test suite (183 tests)
- README with command reference and architecture docs
