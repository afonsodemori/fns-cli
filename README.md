# fns-cli

CLI tools. In progress... Questions? https://afonso.dev/contact

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
