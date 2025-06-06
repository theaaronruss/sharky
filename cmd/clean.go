package cmd

import (
	"fmt"
	"os"
	"sharky/api"
	"sharky/config"
	"strings"

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
			deviceList, err := api.GetDeviceList(accessToken)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			for _, device := range deviceList {
				if strings.EqualFold(device.Name, args[0]) {
					dsn = device.Dsn
				}
			}
			if dsn == "" {
				fmt.Fprintf(os.Stderr, "Failed to find device with name %s\n", args[0])
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
