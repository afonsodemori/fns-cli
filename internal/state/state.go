package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	LatestVersion string    `json:"latest_version"`
	CheckedFor    string    `json:"checked_for"`
	LastCheck     time.Time `json:"last_check"`
}

func GetStateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not get HOME directory: %w", err)
	}

	return filepath.Join(home, ".fns-cli", "state.json"), nil
}

func Load() (*State, error) {
	statePath, err := GetStateFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, fmt.Errorf("Could not read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("Could not parse state file: %w", err)
	}

	return &state, nil
}

func (s *State) Save() error {
	statePath, err := GetStateFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(statePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("Could not create directory: %w", err)
		}
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("Could not marshal state: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0o644); err != nil {
		return fmt.Errorf("Could not write state file: %w", err)
	}

	return nil
}
