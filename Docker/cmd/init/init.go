package init

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var InitCommand = &cobra.Command{
	Use: "init",
	Short: "Init container process run user's process in container. Do not call it outside",
	Run: func(cmd *cobra.Command, args []string)  {
		log.Infof("init come on")
		if len(args) < 1 {
			log.Fatal("missing container command")
		}
		cmdStr := args[0]
		log.Infof("command %s", cmdStr)
		// if err := container.Rim
		// ToDo:
	}

}