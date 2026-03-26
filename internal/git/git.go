package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/afonsodemori/fns-cli/internal/config"
)

type Gist struct {
	ID    string              `json:"id"`
	Files map[string]GistFile `json:"files"`
}

type GistFile struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Language string `json:"language"`
	RawURL   string `json:"raw_url"`
	Size     int    `json:"size"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func GetCurrentBranch() (string, error) {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func GetLocalBranches() ([]string, error) {
	out, err := exec.Command("git", "branch", "--format=%(refname:short)").Output()
	if err != nil {
		return nil, err
	}

	branches := strings.Split(strings.TrimSpace(string(out)), "\n")
	return branches, nil
}

func GetRemoteURL() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func GetGitLabProjectNamespace() (string, error) {
	remoteURL, err := GetRemoteURL()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(?:git@gitlab\.com:|https://gitlab\.com/)(.*)\.git`)
	matches := re.FindStringSubmatch(remoteURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("Can't find namespace for this repository: %s", remoteURL)
	}

	return matches[1], nil
}

// TODO: Probably not here
func ParseIssueKey(param string) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`^((?P<project>\w+)-)?(?P<id>\d+)`)
	matches := re.FindStringSubmatch(param)
	if matches == nil {
		return "", fmt.Errorf("Invalid param \"%s\".", param)
	}

	project := cfg.Jira.DefaultProjectKey
	if p := matches[re.SubexpIndex("project")]; p != "" {
		project = p
	}

	id := matches[re.SubexpIndex("id")]
	if id == "" {
		return "", fmt.Errorf("Invalid param \"%s\".", param)
	}

	return fmt.Sprintf("%s-%s", project, id), nil
}
