package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "sharky",
	Short: "Control your Shark robot vacuum",
	Long: "Sharky is a CLI tool for controlling your Shark robot vacuum",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
