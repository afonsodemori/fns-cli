package git

import (
	"github.com/spf13/cobra"
)

var GitCmd = &cobra.Command{
	Use:   "git",
	Short: "Manage local git repository, and Gitlab operations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
