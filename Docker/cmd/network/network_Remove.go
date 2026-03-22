package network

import (
	"mini-docker/Docker/network"

	"github.com/spf13/cobra"
)

var networkRemoveCmd = &cobra.Command{
	Use:   "commit CONTAINER",
	Short: "commit a container into image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return network.DeleteNetwork(args[0])
	},
}
