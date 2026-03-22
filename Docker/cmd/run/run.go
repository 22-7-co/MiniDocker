package run

import (
	"mini-docker/Docker/Cgroup/subsystem"
	"mini-docker/Docker/runtime"

	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run [flags] [command] [args...]",
	Short: `Create a container with namespace and cgroups limit mydocker run -ti [command]`,
	Long:  `mydocker run -ti --name mycontainer -v /host:/container -p 8080:80 nginx`,
	Example: `  mydocker run -ti --name web -p 8080:80 nginx
  mydocker run -d --mem 512m --cpuset 0-1 busybox sleep 1000`,
	Args: cobra.MinimumNArgs(1), // 至少有一个参数
	RunE: func(cmd *cobra.Command, args []string) error {
		tty, _ := cmd.Flags().GetBool("ti")
		detach, _ := cmd.Flags().GetBool("d")
		mem, _ := cmd.Flags().GetString("mem")
		cpushare, _ := cmd.Flags().GetString("cpushare")
		cpuset, _ := cmd.Flags().GetString("cpuset")
		volume, _ := cmd.Flags().GetString("v")
		name, _ := cmd.Flags().GetString("name")
		envs, _ := cmd.Flags().GetStringSlice("e")
		net, _ := cmd.Flags().GetString("net")
		ports, _ := cmd.Flags().GetStringSlice("p")
		interactive, _ := cmd.Flags().GetBool("interactive")

		resConfig := &subsystem.ResourceConfig{
			MemoryLimit: mem,
			CpuShare:    cpushare,
			CpuSet:      cpuset,
		}

		runtime.Run(tty, detach, args, resConfig, volume, name, envs, net, ports, interactive)
		return nil
	},
}

func init() {
	RunCmd.Flags().BoolP("ti", "t", false, "enable tty")
	RunCmd.Flags().Bool("d", false, "detach container")
	RunCmd.Flags().String("cpushare", "", "cpu share limit")
	RunCmd.Flags().String("mem", "", "memory limit")
	RunCmd.Flags().String("cpuset", "", "cpuset limit")
	RunCmd.Flags().StringP("v", "", "", "volume")
	RunCmd.Flags().String("name", "", "container name")
	RunCmd.Flags().StringSliceP("e", "", nil, "set environment")
	RunCmd.Flags().String("net", "", "container network")
	RunCmd.Flags().StringSliceP("p", "", nil, "port mapping")
	RunCmd.Flags().BoolP("interactive", "i", false, "Keep STDIN open even if not attached")

}
