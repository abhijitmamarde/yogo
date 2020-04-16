package cmd

import (
	"github.com/spf13/cobra"
)

var version = ""

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "App version",
	Run: func(cmd *cobra.Command, args []string) {
		info("v" + version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
