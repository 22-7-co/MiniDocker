package rm

import (
	"mini-docker/Docker/runtime"

	"github.com/spf13/cobra"
)

var RmCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runtime.RemoveContainer(args[0])
	},
}
