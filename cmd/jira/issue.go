package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue [issue-key]",
	Short: "Get Issue information",
	Run: func(cmd *cobra.Command, args []string) {
		issueKey := "N/A"
		if len(args) > 0 {
			issueKey = args[0]
		}
		fmt.Printf("Getting information for issue %s...\n", issueKey)
	},
}

func init() {
	JiraCmd.AddCommand(issueCmd)
	issueCmd.Flags().StringP("comment", "c", "", "Add a comment to the issue")
}
