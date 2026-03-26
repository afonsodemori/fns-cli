package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/afonsodemori/fns-cli/internal/config"
)

type GitLabClient struct {
	cfg *config.Config
}

func NewGitLabClient(cfg *config.Config) *GitLabClient {
	return &GitLabClient{cfg: cfg}
}

func (c *GitLabClient) httpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func (c *GitLabClient) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	u := fmt.Sprintf("%s/%s", c.cfg.GitLab.APIBaseURL, endpoint)
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.GitLab.Token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "fns-cli/0.0.1 (+https://afonso.dev/fns-cli)")

	return req, nil
}

func (c *GitLabClient) GetProjectByNamespace(namespace string) (*Project, error) {
	endpoint := fmt.Sprintf("projects/%s", url.PathEscape(namespace))
	req, err := c.newRequest("GET", endpoint, nil)
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
		return nil, fmt.Errorf("GitLab API error (%d): %s", resp.StatusCode, string(body))
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *GitLabClient) GetPipelines(projectID int, branch string) ([]Pipeline, error) {
	endpoint := fmt.Sprintf("projects/%d/pipelines?ref=%s", projectID, url.QueryEscape(branch))
	req, err := c.newRequest("GET", endpoint, nil)
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
		return nil, fmt.Errorf("GitLab API error (%d): %s", resp.StatusCode, string(body))
	}

	var pipelines []Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	return pipelines, nil
}

func (c *GitLabClient) CreateMergeRequest(namespace string, payload map[string]interface{}) (*MergeRequest, error) {
	endpoint := fmt.Sprintf("projects/%s/merge_requests", url.PathEscape(namespace))
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitLab API error (%d): %s", resp.StatusCode, string(body))
	}

	var mr MergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return nil, err
	}

	return &mr, nil
}

func (c *GitLabClient) GetMergeRequests(namespace, sourceBranch string) ([]MergeRequest, error) {
	endpoint := fmt.Sprintf("projects/%s/merge_requests?source_branch=%s", url.PathEscape(namespace), url.QueryEscape(sourceBranch))
	req, err := c.newRequest("GET", endpoint, nil)
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
		return nil, fmt.Errorf("GitLab API error (%d): %s", resp.StatusCode, string(body))
	}

	var mrs []MergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mrs); err != nil {
		return nil, err
	}

	return mrs, nil
}

type ProjectMap struct {
	Namespaces map[string]int `json:"map"`
}

func (c *GitLabClient) GetCachedProjectID(namespace string) (int, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return 0, err
	}

	path := filepath.Join(home, ".fns-cli", "gitlab-namespace-to-id.map.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	var pm ProjectMap
	if err := json.Unmarshal(data, &pm); err != nil {
		return 0, err
	}

	if pm.Namespaces == nil {
		return 0, nil
	}

	return pm.Namespaces[namespace], nil
}

func (c *GitLabClient) CacheProjectID(namespace string, id int) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(home, ".fns-cli", "gitlab-namespace-to-id.map.json")

	var pm ProjectMap
	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &pm)
	}

	if pm.Namespaces == nil {
		pm.Namespaces = make(map[string]int)
	}

	pm.Namespaces[namespace] = id

	data, err = json.MarshalIndent(pm, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(path, data, 0o644)
}
