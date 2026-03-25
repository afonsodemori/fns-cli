package jira

import (
	"fmt"
	"strings"

	"github.com/afonsodemori/fns-cli/internal/config"
)

type User struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
}

type Issue struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Project     string    `json:"project"`
	Attachments []string  `json:"attachments"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Reporter    User      `json:"reporter"`
	Assignee    *User     `json:"assignee"`
	Sprints     []string  `json:"sprints"`
	TimeSpent   string    `json:"timeSpent"`
	Creator     User      `json:"creator"`
	// Created     time.Time `json:"created"` TODO: Check conversion issue
	// Updated     time.Time `json:"updated"`
}

func GetIssueURL(cfg *config.Config, issueKey string) string {
	return fmt.Sprintf("%s/browse/%s", cfg.Jira.WebBaseURL, strings.ToUpper(issueKey))
}
