package cmd

import (
	"fmt"
	"os"
	"sharky/api"
	"sharky/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use: "clean",
	Short: "Start cleaning home",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accessToken, err := config.GetAccessToken()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
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
		areasToClean := api.Datapoint{
			Name: "SET_Areas_To_Clean",
			Dsn: dsn,
			Value: "*",
		}
		operatingMode := api.Datapoint{
			Name: "SET_Operating_Mode",
			Dsn: dsn,
			Value: fmt.Sprint(api.ModeStart),
		}
		powerMode := api.Datapoint{
			Name: "SET_Power_Mode",
			Dsn: dsn,
			Value: fmt.Sprint(api.PowerNormal),
		}
		err = api.BatchUpdateDatapoints(accessToken, deviceInfo.UserUuid, areasToClean, operatingMode, powerMode)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to start vacuum")
		}
	},
}
