package cmd

import (
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
	rootCmd.AddCommand()
	rootCmd.AddCommand()
}
