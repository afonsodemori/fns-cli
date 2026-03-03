package jira

import (
	"fmt"

	"github.com/spf13/cobra"
)

var transitionCmd = &cobra.Command{
	Use:   "transition [issue-key]",
	Short: "Transition an Issue to another Status",
	Run: func(cmd *cobra.Command, args []string) {
		issueKey := "N/A"
		if len(args) > 0 {
			issueKey = args[0]
		}
		fmt.Printf("Transitioning issue %s...\n", issueKey)
	},
}

func init() {
	JiraCmd.AddCommand(transitionCmd)
}
