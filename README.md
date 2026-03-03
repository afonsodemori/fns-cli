# fns-cli

`fns-cli` is a Go-based CLI tool designed to assist with daily developer tasks, specifically focusing on integrations with **Jira** and **GitLab**. It provides a command-line interface for managing issues, pipelines, and git-related workflows.

> [!NOTE]
> **Work in Progress:** This project is currently being migrated from a legacy version that I developed and use. While the command structures are visible, the underlying implementations are still under development.

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

```bash
go install github.com/afonsodemori/fns-cli@latest
```

## Config

`$HOME/.fns-cli/config.json`

```json
{
  "gitlab": {
    "api_base_url": "https://gitlab.com/api/v4",
    "user_id": 1,
    "token": "aHR0cHM6Ly95b3V0dS5iZS9vYXZNdFVXREJUTQ==",
    "default_project_key": "FCLI"
  },

  "jira": {
    "web_base_url": "https://company.atlassian.net",
    "api_base_url": "https://company.atlassian.net/rest/api/3",
    "token": "aHR0cHM6Ly95b3V0dS5iZS9vYXZNdFVXREJUTQ=="
  }
}
```
