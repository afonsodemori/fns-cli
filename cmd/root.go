package cmd

import (
	"os"

	"github.com/afonsodemori/fns-cli/cmd/config"
	"github.com/afonsodemori/fns-cli/cmd/git"
	"github.com/afonsodemori/fns-cli/cmd/jira"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fns-cli",
	Short: "A CLI to help with daily developer tasks",
	Long:  "A CLI to help with daily developer tasks, including Jira and Gitlab integrations.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(jira.JiraCmd)
	rootCmd.AddCommand(git.GitCmd)
	rootCmd.AddCommand(config.ConfigCmd)
}
