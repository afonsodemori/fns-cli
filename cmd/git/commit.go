package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/afonsodemori/fns-cli/internal/config"
	"github.com/afonsodemori/fns-cli/internal/git"
	"github.com/afonsodemori/fns-cli/internal/jira"
	"github.com/afonsodemori/fns-cli/internal/ui"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:          "commit <message...>",
	Short:        "Commit current staged changes",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			ui.HandleError(err)
		}

		branch, err := git.GetCurrentBranch()
		if err != nil {
			ui.HandleError(err)
		}

		message := strings.Join(args, " ")
		messages := []string{message}

		issueKey, err := git.ParseIssueKey(branch)
		if err == nil {
			issueURL := jira.GetIssueURL(cfg, issueKey)
			messages[0] = fmt.Sprintf("%s: %s", issueKey, messages[0])
			messages = append(messages, fmt.Sprintf("Refs: %s", issueURL))
		}

		commitArgs := []string{"commit"}
		for _, m := range messages {
			commitArgs = append(commitArgs, "-m", m)
		}

		commitCmd := exec.Command("git", commitArgs...)
		commitCmd.Stdout = os.Stdout
		commitCmd.Stderr = os.Stderr
		if err := commitCmd.Run(); err != nil {
			// git commit failure is already printed to Stderr
			os.Exit(1)
		}

		fmt.Println()
		logCmd := exec.Command("git", "log", "-n", "1")
		logCmd.Stdout = os.Stdout
		logCmd.Stderr = os.Stderr
		if err := logCmd.Run(); err != nil {
			ui.HandleError(err)
		}
	},
}

func init() {
	GitCmd.AddCommand(commitCmd)
}
