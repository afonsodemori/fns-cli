package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (c *JiraClient) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
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

	return req, nil
}

func (c *JiraClient) GetIssue(issueKey string) (*Issue, error) {
	url := fmt.Sprintf("%s/issue/%s", c.cfg.Jira.APIBaseURL, issueKey)
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

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

func (c *JiraClient) FindAssignableUsers(issueKey string) ([]User, error) {
	url := fmt.Sprintf("%s/user/assignable/search?issueKey=%s&maxResults=500&recommend=true", c.cfg.Jira.APIBaseURL, issueKey)
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Jira API error (%d): %s", resp.StatusCode, string(body))
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func (c *JiraClient) AssignIssue(issueKey string, user *User) error {
	url := fmt.Sprintf("%s/issue/%s/assignee", c.cfg.Jira.APIBaseURL, issueKey)

	payload := struct {
		AccountID string `json:"accountId"`
	}{
		AccountID: user.AccountID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest("PUT", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Jira API error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *JiraClient) GetTransitions(issueKey string) ([]Transition, error) {
	url := fmt.Sprintf("%s/issue/%s/transitions", c.cfg.Jira.APIBaseURL, issueKey)
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

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
		Transitions []Transition `json:"transitions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return apiResponse.Transitions, nil
}

func (c *JiraClient) DoTransition(issueKey string, transition Transition) error {
	url := fmt.Sprintf("%s/issue/%s/transitions", c.cfg.Jira.APIBaseURL, issueKey)

	payload := struct {
		Transition struct {
			ID string `json:"id"`
		} `json:"transition"`
	}{
		Transition: struct {
			ID string `json:"id"`
		}{
			ID: transition.ID,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Jira API error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}
