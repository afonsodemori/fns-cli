# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`fns-cli` is a Go CLI tool that integrates Jira and GitLab workflows for developer automation. It auto-infers Jira issue keys from git branch names and provides commands to manage issues, commits, merge requests, and pipelines.

Runtime config lives at `~/.fns-cli/config.toml` (see `config.toml.example` for structure).

## Commands

```bash
# Build dev binary for linux/arm64 and symlink to /usr/local/bin/fns-cli
make dev

# Build snapshot with GoReleaser
make build-snapshot

# Test release process without publishing
make release-test

# Start mock API on port 3000 (used during development)
make mock-api
```

There are currently no test files in the project (`go test ./...` returns no test files).

## Architecture

**Entry point:** `main.go` → `cmd.Execute()` (cobra root command in `cmd/root.go`)

**Package layout:**

- `cmd/` — CLI command handlers organized by subcommand group (`jira/`, `git/`, `config/`)
- `internal/config/` — Config loading from `~/.fns-cli/config.json`
- `internal/jira/` — Jira data models (`jira.go`) and HTTP client (`client.go`) using Basic Auth
- `internal/git/` — Git utilities (`git.go`), GitLab models (`gitlab.go`), GitLab API client (`gitlab_client.go`) using Bearer token, GitHub Gist client (`client.go`)
- `internal/ui/` — Terminal output helpers (lipgloss styling, huh interactive forms)
- `internal/state/` — Persistent version-check state at `~/.fns-cli/state.json`
- `internal/version/` — Update checking and self-update logic

**Key patterns:**

1. **Issue key inference** — Most commands accept an optional issue key argument; if omitted, `git.ParseIssueKey()` extracts it from the current git branch name (pattern: `PROJECT-1234`).

2. **GitLab project inference** — The GitLab project namespace is parsed from the git remote URL, so no explicit project config is needed.

3. **Command flow** (e.g., transition): load config → infer issue key from branch → call Jira API → interactive UI selection → call Jira API again → styled output.

4. **Cobra command registration** — Each subgroup has a root command file (e.g., `cmd/jira/jira.go`) that registers sub-commands and is itself added to root in `cmd/root.go`.

## Documentation

A Vitepress documentation project is mounted at `/fns-cli/docs` (symlinked from `/fns-cli-docs/docs/guide/commands/`). When modifying commands in this project, update the corresponding doc files in `docs/`.

## Release

Releases are handled by GoReleaser (`.goreleaser.yaml`). Version is injected via ldflags:

```
-X github.com/afonsodemori/fns-cli/cmd.version={{.Version}}
```

CI publishes on `v*.*.*` tags via `.github/workflows/release.yml`, which also updates a Homebrew tap.
