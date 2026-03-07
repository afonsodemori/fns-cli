<p align="center">
  <img alt="logo: fns-cli" src="https://fns-cli.afonso.dev/assets/logo-banner-transparent.png">
</p>

`fns-cli` is a Go-based CLI tool designed to assist with daily developer tasks, specifically focusing on integrations with **Jira** and **GitLab**. It provides a command-line interface for managing issues, pipelines, and git-related workflows.

> [!NOTE]
> **Work in Progress:** This project is currently being migrated from a legacy version that I developed and use. While the command structures are visible, the underlying implementations are still under development.

Learn more: https://fns-cli.afonso.dev

---

## Main Technologies

- **Go**: Primary programming language (version 1.25.7+).
- **Cobra**: CLI framework for defining commands and flags.
- **Lipgloss**: UI styling library for terminal output.
- **Git**: Direct interaction with local git repositories via shell commands.

## Architecture

- `main.go`: Entry point that executes the root command.
- `cmd/`: Contains the CLI command definitions using Cobra.
  - `root.go`: Root command and subcommand registration.
  - `git/`: Git-specific commands (branch, commit, mr, pipelines).
  - `jira/`: Jira-specific commands (issue, assign, link, transition).
- `internal/`: Core business logic and helpers.
  - `config/`: Handles configuration loading from `$HOME/.fns-cli/config.json`.
  - `git/`: Git-related utilities (e.g., getting current branch).
  - `ui/`: UI helper functions and centralized error handling using Lipgloss.

---

## Installation

```sh
curl -fsSL https://fns-cli.afonso.dev/install.sh | sh
```

or

```sh
go install github.com/afonsodemori/fns-cli@latest
```

Configuration instructions at https://fns-cli.afonso.dev/guide/getting-started.html
