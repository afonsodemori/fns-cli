package git

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/afonsodemori/fns-cli/internal/config"
	"github.com/afonsodemori/fns-cli/internal/git"
	"github.com/afonsodemori/fns-cli/internal/jira"
	"github.com/afonsodemori/fns-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var mergeRequestCmd = &cobra.Command{
	Use:     "merge-request",
	Short:   "Get or Create a Merge Request for the current branch",
	Aliases: []string{"mr"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		sourceBranch, err := git.GetCurrentBranch()
		if err != nil {
			return err
		}

		projectNamespace, err := git.GetGitLabProjectNamespace()
		if err != nil {
			return err
		}

		gitlabClient := git.NewGitLabClient(cfg)
		mrs, err := gitlabClient.GetMergeRequests(projectNamespace, sourceBranch)
		if err != nil {
			return err
		}

		if len(mrs) > 0 {
			for _, mr := range mrs {
				fmt.Println()
				line1 := fmt.Sprintf(" %s ", mr.Title)
				line2 := fmt.Sprintf(" %s ", mr.References.Full)

				var line1Style, line2Style lipgloss.Style
				if mr.State == "opened" {
					line1Style = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("42")).Bold(true)
					line2Style = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
				} else if mr.State == "closed" {
					line1Style = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("196")).Bold(true)
					line2Style = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
				} else {
					line1Style = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("244")).Bold(true)
					line2Style = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
				}

				fmt.Println(line1Style.Render(line1))
				fmt.Println(line2Style.Render(line2))
				fmt.Println()

				mrJSON, _ := json.MarshalIndent(mr, "", "  ")
				fmt.Println(string(mrJSON))
			}

			return printPipelinesInformation()
		}

		confirmCreate, err := ui.Confirm(fmt.Sprintf("No open merge requests for the branch \"%s\". Create one?", sourceBranch))
		if err != nil {
			return err
		}

		if !confirmCreate {
			return printPipelinesInformation()
		}

		title := ""
		re := regexp.MustCompile(`^(?P<id>[A-z0-9]+-[0-9]+)[-_](?P<description>.*)`)
		matches := re.FindStringSubmatch(sourceBranch)
		if len(matches) > 0 {
			id := matches[re.SubexpIndex("id")]
			desc := matches[re.SubexpIndex("description")]
			desc = strings.ReplaceAll(desc, "-", " ")
			desc = strings.ReplaceAll(desc, "_", " ")
			title = fmt.Sprintf("%s: %s", id, desc)
		}

		issueKey, _ := git.ParseIssueKey(sourceBranch)
		jiraLink := jira.GetIssueURL(cfg, issueKey)

		targetBranch, err := git.GetDevelopmentBranch()
		if err != nil {
			return err
		}

		payload := map[string]interface{}{
			"id":            projectNamespace,
			"source_branch": sourceBranch,
			"target_branch": targetBranch,
			"title":         fmt.Sprintf("Draft: %s", title),
			"description":   fmt.Sprintf("Refs: %s", jiraLink),
			"squash":        true,
			"assignee_id":   cfg.GitLab.UserID,
		}

		payloadJSON, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Println(string(payloadJSON))

		confirmContinue, err := ui.Confirm("Continue?")
		if err != nil {
			return err
		}

		if !confirmContinue {
			return nil
		}

		createdMR, err := gitlabClient.CreateMergeRequest(projectNamespace, payload)
		if err != nil {
			return err
		}

		fmt.Println()
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("42")).Bold(true)
		fmt.Println(style.Render(fmt.Sprintf(" %s ", createdMR.Title)))
		style2 := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
		fmt.Println(style2.Render(fmt.Sprintf(" %s ", createdMR.References.Full)))
		fmt.Println()

		createdMRJSON, _ := json.MarshalIndent(createdMR, "", "  ")
		fmt.Println(string(createdMRJSON))

		return printPipelinesInformation()
	},
	SilenceUsage: true,
}

func printPipelinesInformation() error {
	cmd := exec.Command("fns-cli", "git", "pipelines")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	GitCmd.AddCommand(mergeRequestCmd)
}
