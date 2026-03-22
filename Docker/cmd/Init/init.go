package Init

import (
	"mini-docker/Docker/container"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:    "init",
	Short:  "Init container process (internal use)",
	Long:   "Init container process run user's process in container. Do not call it outside.",
	Hidden: true, // 不在 help 中显示
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("init come on")
		return container.RunContainerInitProcess(args[0], args[1:])
	},
}
