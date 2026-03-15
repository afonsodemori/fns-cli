package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/afonsodemori/fns-cli/internal/config"
	"github.com/afonsodemori/fns-cli/internal/git"
	"github.com/afonsodemori/fns-cli/internal/ui"
	"github.com/spf13/cobra"
)

// TODO: Implement "just to work". Review and improve with "git" and "raw" types
var importExtrasCmd = &cobra.Command{
	Use:   "import-extras",
	Short: "Import extra configs (from gist)",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: After install, check if imports are in .bashrc, suggest adding it there, if not present yet.
		cfg, err := config.Load()
		if err != nil {
			ui.HandleError(err)
		}

		client := git.NewClient(cfg)
		gist, err := client.GetGist(cfg.Extras[0].ID)

		home, err := os.UserHomeDir()
		if err != nil {
			// return err -- TODO: Check error.
		}

		extrasDir := filepath.Join(home, ".fns-cli", "extras.d")
		os.RemoveAll(extrasDir)
		os.MkdirAll(extrasDir, 0o755)

		for _, file := range gist.Files {
			// Only process shell scripts
			if file.Type != "application/x-sh" {
				continue
			}

			// Full path for the output file
			filename := fmt.Sprintf("%s-%s", gist.ID, file.Filename)
			outPath := filepath.Join(extrasDir, filename)

			// Write file content
			if err := os.WriteFile(outPath, []byte(file.Content), 0o755); err != nil {
				// return fmt.Errorf("failed to write %s: %w", file.Filename, err) -- TODO: Check error.
			}

			fmt.Printf("Saved: %s\n", outPath)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(importExtrasCmd)
}
