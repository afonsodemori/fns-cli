package cmd

import (
	"fmt"

	"github.com/afonsodemori/fns-cli/internal/state"
	"github.com/afonsodemori/fns-cli/internal/ui"
	ver "github.com/afonsodemori/fns-cli/internal/version"
	"github.com/spf13/cobra"
)

var (
	version = "0.0.0-dev"
	commit  = "none"
	date    = "unknown"
	check   bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of fns-cli",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("fns-cli %s (commit: %s, built: %s)\n", version, commit, date)
		if ver.IsDev(version) {
			ui.Info("You are currently running a development version of fns-cli.")
		}

		// TODO: This is called so state.json is updated when checking
		if check {
			_, _, err := ver.CheckForUpdate(version)
			if err != nil {
				fmt.Printf("Error checking for updates: %v\n", err)
			}
		}

		s, err := state.Load()
		if err == nil && s.LatestVersion != "" {
			if ver.IsNewer(s.LatestVersion, version) {
				fmt.Printf("\nNew version found: %s\n", s.LatestVersion)
				fmt.Println("Run 'fns-cli update' to install it.")
			}

			if check {
				fmt.Printf("\nUpdate state information:\n")
				fmt.Printf("  Latest version: %s\n", s.LatestVersion)
				fmt.Printf("  Checked for:    %s\n", s.CheckedFor)
				fmt.Printf("  Last check:     %s\n", s.LastCheck.Format("2006-01-02 15:04:05"))
			}
		}

		return nil
	},
}

func init() {
	versionCmd.Flags().BoolVarP(&check, "check", "c", false, "Check for new versions")
	rootCmd.AddCommand(versionCmd)
}
