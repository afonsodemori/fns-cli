package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link [issue-key]",
	Short: "Get the URL for an Issue",
	Run: func(cmd *cobra.Command, args []string) {
		issueKey := "N/A"
		if len(args) > 0 {
			issueKey = args[0]
		}
		fmt.Printf("Getting URL for issue %s...\n", issueKey)
	},
}

func init() {
	JiraCmd.AddCommand(linkCmd)
}
