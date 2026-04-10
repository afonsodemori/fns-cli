# fns-cli

A CLI to streamline your daily developer workflow with Jira and GitLab integrations.

## Features

- Auto-infers Jira issue keys from git branch names (e.g. `PROJECT-1234-my-feature`)
- Auto-infers GitLab project namespace from the git remote URL
- Interactive terminal UI for selections and confirmations
- Self-update support

## Installation

### Homebrew (macOS / Linux)

```bash
brew install afonsodemori/tap/fns-cli
```

### Package managers (deb / rpm / apk)

Download the appropriate package from the [releases page](https://github.com/afonsodemori/fns-cli/releases).

## Configuration

Create `~/.fns-cli/config.json`:

```json
{
  "jira": {
    "web_base_url": "https://your-org.atlassian.net",
    "api_base_url": "https://your-org.atlassian.net/rest/api/3",
    "email": "you@example.com",
    "token": "your-jira-api-token",
    "default_project_key": "PROJECT"
  },
  "gitlab": {
    "api_base_url": "https://gitlab.com/api/v4",
    "user_id": 12345,
    "token": "your-gitlab-personal-access-token"
  },
  "extras": [
    {
      "type": "gist",
      "id": "your-gist-id",
      "token": "your-github-token"
    }
  ]
}
```

## Commands

Most commands accept an optional `[issue-key]` argument. If omitted, the issue key is inferred from the current git branch name (e.g. branch `PROJ-42-my-feature` → issue key `PROJ-42`).

### `jira`

```bash
# Show issue details
fns-cli jira issue [issue-key]
fns-cli jira issue [issue-key] --short   # basic info only

# Transition issue to a new status (interactive)
fns-cli jira transition [issue-key]

# Assign issue to a user (interactive)
fns-cli jira assign [issue-key]

# Print the issue URL
fns-cli jira link [issue-key]
```

### `git`

```bash
# List local branches with their GitLab merge request state
fns-cli git branch

# Commit staged changes — auto-prefixes the commit message with the issue key
# and appends a Jira issue reference link
fns-cli git commit <message>

# Show or create a merge request for the current branch (alias: mr)
fns-cli git merge-request
fns-cli git mr

# Show recent pipelines for the current branch
fns-cli git pipelines
```

### `config`

```bash
# Download extra shell scripts from a configured GitHub Gist
fns-cli config import-extras
```

### `version` / `update`

```bash
fns-cli version           # print version info
fns-cli version --check   # also check for a newer release
fns-cli update            # self-update to the latest release
```

## Development

```bash
# Build dev binary for linux/arm64 and symlink to /usr/local/bin/fns-cli
make dev

# Build snapshot with GoReleaser
make build-snapshot

# Start mock API on port 3000
make mock-api
```

## License

[MIT](LICENSE)
