package run

import (
	"fmt"
	"mini-docker/Docker/runtime"

	"github.com/spf13/cobra"
)

var RunCommand = &cobra.Command{
	Use:   "run",
	Short: `Create a container with namespace and cgroups limit mydocker run -ti [command]`,
	Run: func(cmd *cobra.Command, args []string) {
		/*
			1.判断参数是否包含command
			2.获取用户指定的command
			3.调用Run function 去准备启动容器
		*/
		if len(args) < 1 {
			fmt.Println(fmt.Errorf("missing container command"))
			return
		}
		cmdStr := args[0]
		tty, err := cmd.Flags().GetBool("ti") // 获取 -ti 参数值
		if err != nil {
			fmt.Printf("get flag ti error: %v\n", err)
			return
		}
		runtime.Run(tty, cmdStr) // 调用你的业务逻辑
	},
}
