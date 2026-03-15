package jira

import (
	"fmt"

	"github.com/afonsodemori/fns-cli/internal/config"
	"github.com/afonsodemori/fns-cli/internal/git"
	"github.com/afonsodemori/fns-cli/internal/jira"
	"github.com/afonsodemori/fns-cli/internal/ui"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link [issue-key]",
	Short: "Get the URL for an Issue",
	Long:  "The Issue Key. E.g.: FCLI-5712, 5712. If not present, try to infer from current git branch name.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			ui.HandleError(err)
		}

		var issueKey string
		if len(args) > 0 {
			issueKey = args[0]
		} else {
			branch, err := git.GetCurrentBranch()
			if err != nil {
				ui.HandleError(err)
			}
			issueKey = branch
		}

		parsedKey, err := git.ParseIssueKey(issueKey)
		if err != nil {
			ui.HandleError(err)
		}

		fmt.Println(jira.GetIssueURL(cfg, parsedKey))
	},
}

func init() {
	JiraCmd.AddCommand(linkCmd)
}
