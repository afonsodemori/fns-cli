package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afonsodemori/fns-cli/internal/config"
)

type JiraClient struct {
	cfg *config.Config
}

func NewClient(cfg *config.Config) *JiraClient {
	return &JiraClient{cfg: cfg}
}

func (c *JiraClient) httpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func (c *JiraClient) GetIssue(issueKey string) (*Issue, error) {
	url := fmt.Sprintf("%s/issue/%s", c.cfg.Jira.APIBaseURL, issueKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	email := c.cfg.Jira.Email
	apiToken := c.cfg.Jira.Token
	auth := fmt.Sprintf("%s:%s", email, apiToken)
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "fns-cli/0.0.1 (+https://afonso.dev/fns-cli)")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Jira API error (%d): %s", resp.StatusCode, string(body))
	}

	var apiResponse struct {
		ID     string `json:"id"`
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
			// Description string `json:"description"`
			IssueType struct {
				Name string `json:"name"`
			} `json:"issuetype"`
			Project struct {
				Name string `json:"name"`
			} `json:"project"`
			Attachment []struct {
				Filename string `json:"filename"`
			} `json:"attachment"`
			Status struct {
				Name string `json:"name"`
			} `json:"status"`
			Priority struct {
				Name string `json:"name"`
			} `json:"priority"`
			Reporter User  `json:"reporter"`
			Assignee *User `json:"assignee"`
			Sprints  []struct {
				Name string `json:"name"`
			} `json:"customfield_10004"`
			// Comment struct {
			// 	Comments []Comment `json:"comments"`
			// } `json:"comment"`
			TimeTracking struct {
				TimeSpent string `json:"timeSpent"`
			} `json:"timetracking"`
			Creator User `json:"creator"`
			// Created time.Time `json:"created"`
			// Updated time.Time `json:"updated"`
		} `json:"fields"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	issue := &Issue{
		ID:      apiResponse.ID,
		Key:     apiResponse.Key,
		Summary: apiResponse.Fields.Summary,
		// Description: apiResponse.Fields.Description,
		Type:      apiResponse.Fields.IssueType.Name,
		Project:   apiResponse.Fields.Project.Name,
		Status:    apiResponse.Fields.Status.Name,
		Priority:  apiResponse.Fields.Priority.Name,
		Reporter:  apiResponse.Fields.Reporter,
		Assignee:  apiResponse.Fields.Assignee,
		TimeSpent: apiResponse.Fields.TimeTracking.TimeSpent,
		Creator:   apiResponse.Fields.Creator,
		// Created:     apiResponse.Fields.Created,
		// Updated:     apiResponse.Fields.Updated,
		// Comments: apiResponse.Fields.Comment.Comments,
	}

	for _, att := range apiResponse.Fields.Attachment {
		issue.Attachments = append(issue.Attachments, att.Filename)
	}

	for _, sprint := range apiResponse.Fields.Sprints {
		issue.Sprints = append(issue.Sprints, sprint.Name)
	}

	return issue, nil
}
