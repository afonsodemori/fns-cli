package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/afonsodemori/fns-cli/internal/config"
)

func GetCurrentBranch() (string, error) {
	out, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
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
