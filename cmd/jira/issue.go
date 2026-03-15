package jira

import (
	"fmt"
	"strings"

	"github.com/afonsodemori/fns-cli/internal/config"
	"github.com/afonsodemori/fns-cli/internal/git"
	"github.com/afonsodemori/fns-cli/internal/jira"
	"github.com/afonsodemori/fns-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue [issue-key]",
	Short: "Get Issue information",
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

		client := jira.NewClient(cfg)
		issue, err := client.GetIssue(parsedKey)
		if err != nil {
			ui.HandleError(err)
		}

		short, _ := cmd.Flags().GetBool("short")

		displayIssue(cfg, issue, short)
	},
}

func init() {
	JiraCmd.AddCommand(issueCmd)
	issueCmd.Flags().BoolP("short", "s", false, "Show only basic information")
}

func displayIssue(cfg *config.Config, issue *jira.Issue, short bool) {
	statusColors := map[string]string{
		"Backlog":           "244", // Gray
		"Blocked":           "9",   // Red
		"has dependency":    "9",   // Red
		"In Progress":       "12",  // Blue
		"Code Review":       "11",  // Yellow
		"Functional Test":   "13",  // Magenta
		"Prepare Release":   "14",  // Cyan
		"Ready for release": "10",  // Green
		"Done":              "10",  // Green
	}

	colorCode, ok := statusColors[issue.Status]
	if !ok {
		colorCode = "12" // Default Blue
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colorCode)).
		Padding(0, 1)

	statusTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorCode)).
		Italic(true)

	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color(colorCode)).
		Bold(true).
		Padding(0, 1)

	reporterStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorCode))

	fmt.Println()
	fmt.Printf("%s %s\n", statusStyle.Render(fmt.Sprintf("%s | %s", strings.ToUpper(issue.Type), issue.Key)), statusTextStyle.Render(issue.Status))
	fmt.Println(summaryStyle.Render(issue.Summary))

	assigneeName := "Unassigned"
	if issue.Assignee != nil {
		assigneeName = issue.Assignee.DisplayName
	}
	fmt.Println(reporterStyle.Render(fmt.Sprintf(" Reporter: %s -> Assignee: %s ", issue.Reporter.DisplayName, assigneeName)))

	if short {
		return
	}

	fmt.Println()
	fmt.Printf("Priority:    %s\n", issue.Priority)
	fmt.Printf("Project:     %s\n", issue.Project)
	fmt.Printf("Sprints:     %s\n", strings.Join(issue.Sprints, ", "))
	// fmt.Printf("Attachments: %s\n", strings.Join(issue.Attachments, ", "))
	// fmt.Printf("Comments:    %d\n", len(issue.Comments))
	// fmt.Printf("Time Spent:  %s\n", issue.TimeSpent)
	// fmt.Printf("Created:     %s\n", issue.Created.Format("2006-01-02 15:04:05"))
	// fmt.Printf("Updated:     %s\n", issue.Updated.Format("2006-01-02 15:04:05"))

	link := jira.GetIssueURL(cfg, issue.Key)
	linkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Italic(true)
	fmt.Println()
	fmt.Printf("%s %s\n", linkStyle.Render("@see:"), linkStyle.Render(link))
}
