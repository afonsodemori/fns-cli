package git

import (
	"fmt"

	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List branches with additional Gitlab info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing branches with additional Gitlab info...")
	},
}

func init() {
	GitCmd.AddCommand(branchCmd)
}
