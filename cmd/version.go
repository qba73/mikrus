package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.0.0"

// versionCmd represents the logs command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show mikrus client version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mikctl version", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
