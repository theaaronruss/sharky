package cmd

import (
	"fmt"
	"os"
	"sharky/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use: "clean",
	Short: "Clean specific rooms or whole home",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken, err := api.GetAccessToken()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		_, err = api.GetDeviceInfo(accessToken, "")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		// TODO: Set the properties for starting the vacuum (create methods in `api` package)
	},
}
