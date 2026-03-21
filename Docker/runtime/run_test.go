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
		{"memory test", args{"memory"}, "/sys/fs/cgroup/memory"},
		{"cpushare test", args{"cpu"}, "/sys/fs/cgroup/cpu,cpuacct"},
		{"cpuset test", args{"cpuset"}, "/sys/fs/cgroup/cpuset"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := subsystem.FindCgroupMountPoint(tt.args.subsystem); err != nil || got != tt.want {
				t.Errorf("FindCgroupMountPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
