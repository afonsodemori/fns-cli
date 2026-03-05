package git

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/afonsodemori/fns-cli/internal/config"
)

type GitClient struct {
	cfg *config.Config
}

func NewClient(cfg *config.Config) *GitClient {
	return &GitClient{cfg: cfg}
}

func (c *GitClient) httpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func (c *GitClient) GetGist(id string) (*Gist, error) {
	url := fmt.Sprintf("https://api.github.com/gists/%s", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.Extras[0].Token))
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
		return nil, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(body))
	}

	var apiResponse *Gist

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return apiResponse, nil
}
