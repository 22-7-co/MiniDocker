package network

import (
	"mini-docker/Docker/network"

	"github.com/spf13/cobra"
)

var networkCreateCmd = &cobra.Command{
	Use:   "commit CONTAINER",
	Short: "commit a container into image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return network.CreateNetwork(args[0], args[1], args[2])
	},
}
