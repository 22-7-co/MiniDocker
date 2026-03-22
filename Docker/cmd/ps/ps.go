package ps

import (
	"mini-docker/Docker/runtime"

	"github.com/spf13/cobra"
)

var PsCmd = &cobra.Command{
	Use:   "ps",
	Short: "show all containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runtime.ListContainers()
	},
}
