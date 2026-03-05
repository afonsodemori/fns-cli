package config

import (
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage fns-cli config",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
