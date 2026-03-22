package network

import "github.com/spf13/cobra"

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "container network commands",
}

func init() {
	// 在 root.go 里已经 AddCommand(networkCmd) 了，这里只加子命令
	NetworkCmd.AddCommand(
		networkCreateCmd,
		networkListCmd,
		networkRemoveCmd,
	)
}
