package jira

import (
	"github.com/spf13/cobra"
)

var JiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Manage Issues",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
