package cmd

import (
	"fmt"
	"os"
	"sharky/config"

	"github.com/spf13/cobra"
)

var remove string
var defaultDsn string

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&remove, "remove", "", "", "Remove config property")
	configCmd.Flags().StringVarP(&defaultDsn, "default-dsn", "", "", "Default DSN to use for commands")
}

var configCmd = &cobra.Command{
	Use: "config",
	Short: "Set and retrieve configuration properties",
	Run: func(cmd *cobra.Command, args []string) {
		if remove != "" {
			err := config.RemoveProperty(remove)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			fmt.Printf("Removed \"%s\" property\n", remove)
		}
		if defaultDsn != "" {
			err := config.SetProperty(config.DefaultDsn, defaultDsn)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			fmt.Printf("Updated default DSN to \"%s\"\n", defaultDsn)
		}
	},
}
