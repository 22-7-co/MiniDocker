package cmd

import (
	"mini-docker/Docker/cmd/Init"
	"mini-docker/Docker/cmd/commit"
	"mini-docker/Docker/cmd/exec"
	"mini-docker/Docker/cmd/log"
	"mini-docker/Docker/cmd/network"
	"mini-docker/Docker/cmd/ps"
	"mini-docker/Docker/cmd/rm"
	"mini-docker/Docker/cmd/run"
	"mini-docker/Docker/cmd/stop"

	"github.com/spf13/cobra"
)

/*
整合所有子命令
*/

var rootCmd = &cobra.Command{
	Use:   "mydocker",
	Short: "迷你 Docker 实现",
	Long:  ``,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {

	rootCmd.AddCommand(
		Init.InitCmd,
		run.RunCmd,
		log.LogCmd,
		stop.StopCmd,
		commit.CommitCmd,
		rm.RmCmd,
		ps.PsCmd,
		exec.ExecCmd,
		network.NetworkCmd,
	)
}
