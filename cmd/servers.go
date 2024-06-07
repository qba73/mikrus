package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// serversCmd represents the servers command
var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "show servers associated with the Mikrus account",
	Long:  `show servers associated with the Mikrus account`,
	Run: func(cmd *cobra.Command, args []string) {
		servers, err := client.Servers()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(servers)
	},
}

func init() {
	rootCmd.AddCommand(serversCmd)
}
