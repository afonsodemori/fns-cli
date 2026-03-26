package git

import (
	"fmt"
	"math"
	"time"

	"github.com/afonsodemori/fns-cli/internal/config"
	"github.com/afonsodemori/fns-cli/internal/git"
	"github.com/afonsodemori/fns-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "Get pipelines for the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			ui.HandleError(err)
		}

		branch, err := git.GetCurrentBranch()
		if err != nil {
			ui.HandleError(err)
		}

		namespace, err := git.GetGitLabProjectNamespace()
		if err != nil {
			ui.HandleError(err)
		}

		client := git.NewGitLabClient(cfg)
		projectID, err := client.GetCachedProjectID(namespace)
		if err != nil {
			ui.HandleError(err)
		}

		if projectID == 0 {
			project, err := client.GetProjectByNamespace(namespace)
			if err != nil {
				ui.HandleError(err)
			}
			projectID = project.ID
			err = client.CacheProjectID(namespace, projectID)
			if err != nil {
				ui.HandleError(err)
			}
		}

		pipelines, err := client.GetPipelines(projectID, branch)
		if err != nil {
			ui.HandleError(err)
		}

		if len(pipelines) == 0 {
			fmt.Println("\nNo pipelines found for this branch.")
			return
		}

		fmt.Println("\nLast pipelines:")
		limit := 10
		if len(pipelines) < limit {
			limit = len(pipelines)
		}

		for _, p := range pipelines[:limit] {
			statusStyle := getStatusStyle(p.Status)
			relativeTime := formatRelativeTime(p.UpdatedAt)

			// TODO: Duration doesn't come from API response. This was the Workaound in the legacy code
			// const pipelineDuration = parseInt(pipeline.duration.split(' ')[0]);
			// if (pipelineDuration > 0 && pipelineDuration < 60) fns += ` (${pipeline.duration})`;
			line := fmt.Sprintf("%s %s", p.Status, relativeTime)
			if p.Duration > 0 && p.Duration < 60 {
				line += fmt.Sprintf(" (%d sec)", p.Duration)
			}

			fmt.Printf(" - %s\n", statusStyle.Render(line))
		}
	},
}

func getStatusStyle(status string) lipgloss.Style {
	base := lipgloss.NewStyle()
	switch status {
	case "running":
		return base.Foreground(lipgloss.Color("12")) // Blue
	case "success":
		return base.Foreground(lipgloss.Color("10")) // Green
	case "failed":
		return base.Foreground(lipgloss.Color("9")) // Red
	case "canceled":
		return base.Foreground(lipgloss.Color("8")) // Gray
	default:
		return base.Foreground(lipgloss.Color("11")) // Yellow
	}
}

func formatRelativeTime(t time.Time) string {
	diff := time.Since(t).Seconds()
	if diff < 60 {
		return fmt.Sprintf("%d sec ago", int(diff))
	} else if diff < 3600 {
		return fmt.Sprintf("%d min ago", int(math.Round(diff/60)))
	} else if diff < 86400 {
		return fmt.Sprintf("%d hr ago", int(math.Round(diff/3600)))
	} else {
		return fmt.Sprintf("%d days ago", int(math.Round(diff/86400)))
	}
}

func init() {
	GitCmd.AddCommand(pipelinesCmd)
}
