package log

import (
	"mini-docker/Docker/runtime"

	"github.com/spf13/cobra"
)

var LogCmd = &cobra.Command{
	Use:   "log container",
	Short: "show container log",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runtime.LogContainer(args[0])
	},
}
