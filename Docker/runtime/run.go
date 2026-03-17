package runtime

import (
	"mini-docker/Docker/container"
	"os"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
		return
	}
	log.Infof("parent process run")
	parent.Wait()
	os.Exit(-1)
}
