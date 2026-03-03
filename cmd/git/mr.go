package git

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mrCmd = &cobra.Command{
	Use:   "mr",
	Short: "Get or Create a Merge Request for the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting or Creating a Merge Request for the current branch...")
	},
}

func init() {
	GitCmd.AddCommand(mrCmd)
}
