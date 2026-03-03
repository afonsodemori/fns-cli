package git

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "Get pipelines for the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting pipelines for the current branch...")
	},
}

func init() {
	GitCmd.AddCommand(pipelinesCmd)
}
