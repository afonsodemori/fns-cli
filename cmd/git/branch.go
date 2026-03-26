package git

import (
	"fmt"
	"strings"

	"github.com/afonsodemori/fns-cli/internal/config"
	"github.com/afonsodemori/fns-cli/internal/git"
	"github.com/afonsodemori/fns-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:          "branch",
	Short:        "List branches with additional Gitlab info",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			ui.HandleError(err)
		}

		gitlabClient := git.NewGitLabClient(cfg)
		namespace, err := git.GetGitLabProjectNamespace()
		if err != nil {
			ui.HandleError(err)
		}

		branches, err := git.GetLocalBranches()
		if err != nil {
			ui.HandleError(err)
		}

		stateStyles := map[string]lipgloss.Style{
			"opened": lipgloss.NewStyle().Foreground(lipgloss.Color("42")),  // green
			"merged": lipgloss.NewStyle().Foreground(lipgloss.Color("33")),  // blue
			"closed": lipgloss.NewStyle().Foreground(lipgloss.Color("196")), // red
			" none ": lipgloss.NewStyle().Foreground(lipgloss.Color("244")), // gray
		}

		for _, branch := range branches {
			branch = strings.TrimSpace(branch)
			if branch == "" {
				continue
			}

			mrs, err := gitlabClient.GetMergeRequests(namespace, branch)
			if err != nil {
				ui.HandleError(err)
			}

			state := " none "
			if len(mrs) > 0 {
				state = mrs[len(mrs)-1].State
			}

			style, ok := stateStyles[state]
			if !ok {
				style = stateStyles[" none "]
			}

			fmt.Printf("%s %s\n", style.Render(fmt.Sprintf("(%s)", state)), branch)
		}
	},
}

func init() {
	GitCmd.AddCommand(branchCmd)
}
