package cmd

import (
	"fmt"
	"os"

	"github.com/qba73/mikrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mikctl",
	Short: "mikctl is a client for Mikrus API",
	Long: `mikctl is a command-line client for the Mikrus VPS hosting service.
It allows you to search for existing servers, logs, databases,
monitor heath parameters, resource usage and execute commands remotely.
You can also inspect your account details and any server details
you have provisioned.

For more information, see https://github.com/qba73/mikrus`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	apiKey string
	srvID  string
	client mikrus.Client
)

func init() {
	viper.SetConfigName(".mikrus")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		os.Exit(1)
	}
	viper.SetEnvPrefix("mikrus")
	viper.AutomaticEnv()
	cobra.OnInitialize(func() {
		client = mikrus.New(viper.GetString("apiKey"), viper.GetString("srvID"))
	})

	rootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "Mikrus server API key")
	viper.BindPFlag("apiKey", rootCmd.PersistentFlags().Lookup("apiKey"))
	viper.BindEnv("apiKey", "MIKRUS_API_KEY")

	rootCmd.PersistentFlags().StringVar(&srvID, "srvID", "", "Mikrus server ID")
	viper.BindPFlag("srvID", rootCmd.PersistentFlags().Lookup("srvID"))
	viper.BindEnv("apiKey", "MIKRUS_SRV_ID")
}
