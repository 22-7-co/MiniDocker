package network

import (
	"mini-docker/Docker/network"

	"github.com/spf13/cobra"
)

var networkListCmd = &cobra.Command{
	Use:   "list ",
	Short: "list all networks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := network.Init(); err != nil {
			return err
		}
		return network.ListNetwork()
	},
}
