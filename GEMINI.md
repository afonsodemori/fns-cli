# fns-cli

## Project Overview

`fns-cli` is a Go-based CLI tool designed to assist with daily developer tasks, specifically focusing on integrations with **Jira** and **GitLab**. It provides a command-line interface for managing issues, pipelines, and git-related workflows.

### Main Technologies

- **Go**: Primary programming language (version 1.25.7+).
- **Cobra**: CLI framework for defining commands and flags.
- **Lipgloss**: UI styling library for terminal output.
- **Git**: Direct interaction with local git repositories via shell commands.

### Architecture

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

## Building and Running

### Build Commands

The project includes a `Makefile` for local installation and completion setup.

- **Build and Install:**

  ```bash
  make install
  ```

  This builds the binary to `bin/fns-cli` and generates bash completions in `$HOME/.fns-cli/bash-completion.sh`.

- **Direct Build:**
  ```bash
  go build -o bin/fns-cli
  ```

### Configuration

The application expects a configuration file at `$HOME/.fns-cli/config.json`.

**Example Configuration:**

```json
{
  "gitlab": {
    "api_base_url": "https://gitlab.com/api/v4",
    "user_id": 1,
    "token": "YOUR_GITLAB_TOKEN"
  },
  "jira": {
    "web_base_url": "https://company.atlassian.net",
    "api_base_url": "https://company.atlassian.net/rest/api/3",
    "token": "YOUR_JIRA_TOKEN",
    "default_project_key": "PROJ"
  }
}
```

---

## Development Conventions

### Command Structure

New commands should be added under the `cmd/` directory and registered in their respective parent command's `init()` function or in `cmd/root.go`.

### Error Handling

Centralized error handling is provided in `internal/ui/output.go`. Use `ui.HandleError(err)` to display a styled error message and exit the application.

### UI Styling

Use **Lipgloss** for all terminal output styling to maintain consistency. Styles should be defined or used within the `internal/ui` package or locally in commands for specific needs.

### Configuration Access

Configuration should be loaded via `config.Load()`, which returns a `*Config` struct. It handles home directory expansion and JSON unmarshaling.

### Testing

- **TODO**: No explicit test suite found. New features should ideally include unit tests for logic in `internal/`.
