package cmd

import (
	"context"
	"fmt"

	"github.com/afonsodemori/fns-cli/internal/ui"
	ver "github.com/afonsodemori/fns-cli/internal/version"
	"github.com/charmbracelet/huh"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:          "update",
	Aliases:      []string{"self-update"},
	Short:        "Updates fns-cli to the latest version",
	SilenceUsage: true, // TODO: Check on other commands
	RunE: func(cmd *cobra.Command, args []string) error {
		if ver.IsDev(version) {
			ui.Warn("You are currently running a development version of fns-cli (installed via 'go install' or 'go build')")

			options := []huh.Option[string]{
				huh.NewOption("Yes, update", "yes"),
				huh.NewOption("No, keep current dev version", "no").Selected(true),
			}
			confirm, _ := ui.Select("Replace dev version?", options)

			if confirm == "no" {
				fmt.Println("😌 Update cancelled.")
				return nil
			}
		}

		latest, isNewer, err := ver.CheckForUpdate(version)

		if !isNewer {
			fmt.Printf("👍 Current version (%s) is already the latest.\n", version)
			return nil
		}

		fmt.Printf("New version found: %s. Downloading and installing...\n", latest)
		// TODO: Updater instantiated here and in internal/version/version.go -- could reuse?
		updater, err := selfupdate.NewUpdater(selfupdate.Config{})
		repo := selfupdate.ParseSlug("afonsodemori/fns-cli")
		_, err = updater.UpdateSelf(context.Background(), version, repo)
		if err != nil {
			return fmt.Errorf("Update failed: %w", err)
		}

		fmt.Println("🎉 Update successful!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
