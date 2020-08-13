package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/caquillo07/rotom-bot/metrics"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version, branch, and commit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(
			"Version: %v\nBranch: %v\nCommit: %v\nBuilt On: %v\n",
			metrics.Version, metrics.Branch, metrics.Commit, metrics.Date,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
