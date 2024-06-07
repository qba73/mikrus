package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "logs lists log entries for the server",
	Long: `Logs lists last 10 log entries for the server
assiociated with the API key and server name`,
	Run: func(cmd *cobra.Command, args []string) {
		logs, err := client.Logs()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(logs)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
