package stop

import (
	"mini-docker/Docker/runtime"

	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use:   "commit CONTAINER",
	Short: "commit a container into image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runtime.StopContainer(args[0])
	},
}
