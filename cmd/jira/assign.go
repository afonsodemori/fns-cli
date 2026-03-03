package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var assignCmd = &cobra.Command{
	Use:   "assign [issue-key]",
	Short: "Assign an Issue to another User",
	Run: func(cmd *cobra.Command, args []string) {
		issueKey := "N/A"
		if len(args) > 0 {
			issueKey = args[0]
		}
		fmt.Printf("Assigning issue %s...\n", issueKey)
	},
}

func init() {
	JiraCmd.AddCommand(assignCmd)
}
