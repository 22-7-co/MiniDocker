package ps

import (
	"mini-docker/Docker/runtime"

	"github.com/spf13/cobra"
)

var PsCmd = &cobra.Command{
	Use:   "commit CONTAINER",
	Short: "commit a container into image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runtime.ListContainers()
	},
}
