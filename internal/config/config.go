package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Jira struct {
		WebBaseURL        string `json:"web_base_url"`
		APIBaseURL        string `json:"api_base_url"`
		Token             string `json:"token"`
		DefaultProjectKey string `json:"default_project_key"`
	} `json:"jira"`

	GitLab struct {
		APIBaseURL string `json:"api_base_url"`
		UserID     int    `json:"user_id"`
		Token      string `json:"token"`
	} `json:"gitlab"`

	Extras []struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Token string `json:"token"`
	} `json:"extras"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Could not get HOME directory: %w", err)
	}

	configPath := filepath.Join(home, ".fns-cli", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Config file not found at %s.", configPath)
		}
		return nil, fmt.Errorf("Could not read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Could not parse config file: %w", err)
	}

	return &config, nil
}
