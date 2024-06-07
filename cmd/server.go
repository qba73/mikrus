package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "show server details",
	Long:  `Show server details.`,
	Run: func(cmd *cobra.Command, args []string) {
		server, err := client.Info()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(server)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
