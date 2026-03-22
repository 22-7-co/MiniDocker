package network

import (
	"mini-docker/Docker/network"

	"github.com/spf13/cobra"
)

var networkCreateCmd = &cobra.Command{
	Use:   "create  [network-name]",
	Short: "create a network",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		driver, _ := cmd.Flags().GetString("driver")
		subnet, _ := cmd.Flags().GetString("subnet")
		name := args[0]

		if err := network.Init(); err != nil {
			return err
		}
		return network.CreateNetwork(driver, subnet, name)
	},
}

func init() {
	networkCreateCmd.Flags().String("driver", "bridge", "network dirver(e.g. bridge, macvlan)")
	networkCreateCmd.Flags().String("subnet", "172.18.0.0/16", "subnet in CIDR format (e.g. 172.20.0.0/16)")
}
