package cmd

import (
	"fmt"
	"sharky/config"

	"github.com/spf13/cobra"
)

var defaultDsn string

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&defaultDsn, "default-dsn", "", "", "Default DSN to use for commands")
}

var configCmd = &cobra.Command{
	Use: "config",
	Short: "Set and retrieve configuration properties",
	Run: func(cmd *cobra.Command, args []string) {
		if defaultDsn != "" {
			config.SetProperty(config.DefaultDsn, defaultDsn)
			fmt.Printf("Updated default DSN to %s\n", defaultDsn)
		}
	},
}
