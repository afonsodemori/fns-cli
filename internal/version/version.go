package version

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/afonsodemori/fns-cli/internal/state"
	"github.com/creativeprojects/go-selfupdate"
)

func IsDev(v string) bool {
	return v == "0.0.0-dev" || v == "" || strings.Contains(v, "SNAPSHOT")
}

func IsNewer(latestVersion, currentVersion string) bool {
	if IsDev(currentVersion) {
		return true
	}

	l, err := semver.NewVersion(latestVersion)
	if err != nil {
		return false
	}
	c, err := semver.NewVersion(currentVersion)
	if err != nil {
		// If current version is not semver but latest is, we consider it newer
		return true
	}

	return l.GreaterThan(c)
}

func CheckForUpdate(currentVersion string) (string, bool, error) {
	updater, err := selfupdate.NewUpdater(selfupdate.Config{})
	if err != nil {
		return "", false, fmt.Errorf("Failed to create updater: %w", err)
	}

	repo := selfupdate.ParseSlug("afonsodemori/fns-cli")
	latest, found, err := updater.DetectLatest(context.Background(), repo)
	if err != nil {
		return "", false, fmt.Errorf("Failed to query GitHub: %w", err)
	}

	if !found {
		return "", false, nil
	}

	s, err := state.Load()
	if err == nil {
		s.LatestVersion = latest.Version()
		s.CheckedFor = currentVersion
		s.LastCheck = time.Now()
		_ = s.Save()
	}

	if IsNewer(latest.Version(), currentVersion) {
		return latest.Version(), true, nil
	}

	return latest.Version(), false, nil
}
