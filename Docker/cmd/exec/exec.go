package exec

import (
	"mini-docker/Docker/runtime"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ExecCmd = &cobra.Command{
	Use:   "commit CONTAINER",
	Short: "commit a container into image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv(runtime.EnvExecPid) != "" {
			log.Infof("pid callback pid %d", os.Getgid())
			return
		}

		containerName := args[0]
		commandArray := args[1:]

		if err := runtime.ExecContainer(containerName, commandArray); err != nil {
			log.Errorf("exec container %s err: %v", containerName, err)
			os.Exit(1)
		}
	},
}
