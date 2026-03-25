package jira

import (
	"fmt"

	"github.com/afonsodemori/fns-cli/internal/config"
	"github.com/afonsodemori/fns-cli/internal/git"
	"github.com/afonsodemori/fns-cli/internal/jira"
	"github.com/afonsodemori/fns-cli/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var transitionCmd = &cobra.Command{
	Use:          "transition [issue-key]",
	Short:        "Transition an Issue to another Status",
	Long:         "The Issue Key. E.g.: FCLI-5712, 5712. If not present, try to infer from current git branch name.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		var issueKey string
		if len(args) > 0 {
			issueKey = args[0]
		} else {
			branch, err := git.GetCurrentBranch()
			if err != nil {
				return err
			}
			issueKey = branch
		}

		parsedKey, err := git.ParseIssueKey(issueKey)
		if err != nil {
			return err
		}

		client := jira.NewClient(cfg)

		issue, err := client.GetIssue(parsedKey)
		if err != nil {
			return err
		}
		displayIssue(cfg, issue, true)

		fmt.Printf("\nFetching transitions for %s...\n", parsedKey)
		transitions, err := client.GetTransitions(parsedKey)
		if err != nil {
			return err
		}

		if len(transitions) == 0 {
			return fmt.Errorf("no transitions found for issue %s", parsedKey)
		}

		options := make([]huh.Option[jira.Transition], len(transitions))
		for i := range transitions {
			options[i] = huh.NewOption(transitions[i].Name, transitions[i])
		}

		selectedTransition, err := ui.Select("Choose a new Status:", options)
		if err != nil {
			return err
		}

		fmt.Printf("Transitioning %s to \"%s\"...\n", parsedKey, selectedTransition.Name)
		err = client.DoTransition(parsedKey, selectedTransition)
		if err != nil {
			return err
		}

		issue, err = client.GetIssue(parsedKey)
		if err != nil {
			return err
		}
		displayIssue(cfg, issue, true)

		reassign, err := ui.Confirm("Reassign Issue?")
		if err != nil {
			return err
		}

		if reassign {
			fmt.Printf("\nFetching assignable users for %s...\n", parsedKey)
			users, err := client.FindAssignableUsers(parsedKey)
			if err != nil {
				return err
			}

			if len(users) == 0 {
				return fmt.Errorf("no assignable users found for issue %s", parsedKey)
			}

			userOptions := make([]huh.Option[*jira.User], len(users))
			for i := range users {
				userOptions[i] = huh.NewOption(users[i].DisplayName, &users[i])
			}

			selectedUser, err := ui.Select("Assign to:", userOptions)
			if err != nil {
				return err
			}

			fmt.Printf("Assigning %s to %s...\n", parsedKey, selectedUser.DisplayName)
			err = client.AssignIssue(parsedKey, selectedUser)
			if err != nil {
				return err
			}

			issue, err = client.GetIssue(parsedKey)
			if err != nil {
				return err
			}
			displayIssue(cfg, issue, true)
		}

		return nil
	},
}

func init() {
	JiraCmd.AddCommand(transitionCmd)
}
