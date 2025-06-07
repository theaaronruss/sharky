package cmd

import (
	"fmt"
	"os"
	"sharky/api"
	"sharky/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use: "stop",
	Short: "Stop the vacuum",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken, err := config.GetAccessToken()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		dsn := ""
		if len(args) > 0 {
			dsn, err = api.GetDeviceDsn(accessToken, args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		} else {
			dsn, err = config.GetProperty(config.DefaultDsn)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Please provide a device name or configure a default DSN")
				return
			}
		}
		deviceInfo, err := api.GetDeviceInfo(accessToken, dsn)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		err = api.UpdateDatapoint(accessToken, deviceInfo.UserUuid, dsn, "SET_Operating_Mode", fmt.Sprint(api.ModeReturn))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}
