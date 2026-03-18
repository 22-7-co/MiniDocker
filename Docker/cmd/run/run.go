package run

import (
	"log"
	"mini-docker/Docker/runtime"

	"github.com/spf13/cobra"
)

var (
	flagTty      bool
	flagMemory   string
	flagCpuset   string
	flagCpuShare string
	flagVolume   string // -v 通常只允許一次，如需多次請改 StringSlice
	flagDetach   bool
	flagName     string
	flagEnvs     []string // -e 可重复
	flagNetwork  string
	flagPorts    []string // -p 可重复
)

var RunCommand = &cobra.Command{
	Use:   "run [command] [args...]",
	Short: `Create a container with namespace and cgroups limit mydocker run -ti [command]`,
	Long:  `mydocker run -ti --name mycontainer -v /host:/container -p 8080:80 nginx`,
	Example: `  mydocker run -ti --name web -p 8080:80 nginx
  mydocker run -d --mem 512m --cpuset 0-1 busybox sleep 1000`,
	Args: cobra.MinimumNArgs(1), // 至少有一个参数
	Run: func(cmd *cobra.Command, args []string) {
		cmdStr := args[0]
		cmdArgs := args[1:]

		tty := flagTty
		memLimit := flagMemory
		cpuset := flagCpuset
		cpuShare := flagCpuShare
		volume := flagVolume
		detach := flagDetach
		name := flagName
		envs := flagEnvs
		net := flagNetwork
		ports := flagPorts

		if err := runtime.Run(tty,
			cmdStr,
			cmdArgs,
			memLimit,
			cpuset,
			cpuShare,
			volume,
			detach,
			name,
			envs,
			net,
			ports,
		); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// -ti / --ti
	RunCommand.Flags().BoolVarP(&flagTty, "ti", "t", false, "enable tty (allocate a pseudo-TTY)")

	// --mem
	RunCommand.Flags().StringVar(&flagMemory, "mem", "", "memory limit (e.g. 512m, 1g)")

	// --cpuset
	RunCommand.Flags().StringVar(&flagCpuset, "cpuset", "", "cpuset CPUs (e.g. 0-3,0,2)")

	// --cpushare
	RunCommand.Flags().StringVar(&flagCpuShare, "cpushare", "", "cpu shares (relative weight)")

	// -v / --v   （這裡只允許一次，如需 -v 多次請改成 StringSliceVar）
	RunCommand.Flags().StringVarP(&flagVolume, "v", "v", "", "bind mount (e.g. /host/path:/container/path)")

	// -d / --detach
	RunCommand.Flags().BoolVarP(&flagDetach, "detach", "d", false, "run container in background and print container ID")

	// --name
	RunCommand.Flags().StringVar(&flagName, "name", "", "assign a name to the container")

	// -e / --env   （可重複）
	RunCommand.Flags().StringSliceVarP(&flagEnvs, "e", "e", nil, "set environment variables (can be repeated)")

	// --net
	RunCommand.Flags().StringVar(&flagNetwork, "net", "default", "container network (bridge, host, none, ...)")

	// -p / --publish   （可重複）
	RunCommand.Flags().StringSliceVarP(&flagPorts, "p", "p", nil, "publish a container's port(s) to the host (can be repeated)")

}
