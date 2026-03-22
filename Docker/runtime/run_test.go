package runtime

import (
	"mini-docker/Docker/Cgroup/subsystem"
	"testing"
)

func TestFindCgroupMountPoint(t *testing.T) {
	type args struct {
		subsystem string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		//v2
		{"memory test", args{"memory"}, "/sys/fs/cgroup"},
		{"cpushare test", args{"cpu"}, "/sys/fs/cgroup"},
		{"cpuset test", args{"cpuset"}, "/sys/fs/cgroup"},
		// v1
		// {"memory test", args{"memory"}, "/sys/fs/cgroup/memory"},
		// {"cpushare test", args{"cpu"}, "/sys/fs/cgroup/cpu"},
		// {"cpuset test", args{"cpuset"}, "/sys/fs/cgroup/cpuset"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := subsystem.FindCgroupMountPoint(tt.args.subsystem)
			if err != nil || got != tt.want {
				t.Errorf("FindCgroupMountPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
在现代 Linux（2022 年后的大多数发行版，如 Ubuntu 22.04+、Fedora、Debian 12+、很多云服务器镜像）
默认都启用 cgroup v2，挂载点统一变成：/sys/fs/cgroup （类型 cgroup2）
里面不再有 memory、cpu、cpuset 等子目录
所有控制器都在同一个树里，通过文件如 cpu.max、memory.max 来控制
*/
