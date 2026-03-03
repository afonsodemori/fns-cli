package git

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit <message...>",
	Short: "Commit current staged changes",
	Run: func(cmd *cobra.Command, args []string) {
		message := strings.Join(args, " ")
		fmt.Printf("Committing current staged changes with message: %s...\n", message)
	},
}

func init() {
	GitCmd.AddCommand(commitCmd)
}
