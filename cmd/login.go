package cmd

import (
	"fmt"
	"sharky/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use: "login",
	Short: "Generate URL for logging into account",
	Run: func(cmd *cobra.Command, args []string) {
		url, codeVerifier := api.GenerateLoginUrl()
		fmt.Println("Code verifier:", codeVerifier)
		fmt.Println()
		fmt.Println(url)
	},
}
