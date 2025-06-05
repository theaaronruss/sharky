package cmd

import (
	"fmt"
	"os"
	"sharky/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(devicesCmd)
}

var devicesCmd = &cobra.Command{
	Use: "devices",
	Short: "Get list of devices attached to your account",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken, err := api.GetAccessToken()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		devices, err := api.GetDeviceList(accessToken)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		for _, device := range devices {
			fmt.Println("*", device.Name, "-", device.Status)
		}
	},
}
