package cmd

import (
	"fmt"
	"os"
	"sharky/config"

	"github.com/spf13/cobra"
)

var list bool
var remove string
var defaultDsn string

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVarP(&list, "list", "", false, "List of config properties")
	configCmd.Flags().StringVarP(&remove, "remove", "", "", "Remove config property")
	configCmd.Flags().StringVarP(&defaultDsn, "default-dsn", "", "", "Default DSN to use for commands")
}

var configCmd = &cobra.Command{
	Use: "config",
	Short: "Set and retrieve configuration properties",
	Run: func(cmd *cobra.Command, args []string) {
		if list {
			// TODO: List all config properties
			properties, err := config.GetAllProperties()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			for _, property := range properties {
				fmt.Println(property)
			}
			return
		}
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
