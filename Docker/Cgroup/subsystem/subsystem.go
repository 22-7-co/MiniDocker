package subsystem

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

/*
Name： 名称，如memory、cpuset、cpushare
Set: 写入配置文件，对资源进行限制
Apply: 将PID加入当前Cgroup
Remove: 将PID移出当前Cgroup
*/
type Subsystem interface {
	Name() string
	Set(containerName string, res *ResourceConfig) error
	Apply(containerName string, pid int) error
	Remove(containerName string) error
}

var Subsystems = make(map[string]Subsystem)
