package subsystem

import (
	"fmt"
	"os"
	"path"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type MemorySubsystem struct{}

func (m MemorySubsystem) Name() string {
	return "memory"
}

func (m MemorySubsystem) Set(cgroupName string, res *ResourceConfig) error {
	// 获取自定义Cgroup的路径，没有则创建，如：/sys/fs/cgroup/memory/mydocker-cgroup
	cgroupPath, err := GetCgroupPath(m.Name(), cgroupName)
	if err != nil {
		return err
	}
	log.Infof("%s cgroup path: %s", m.Name(), cgroupName)
	if err != nil {
		return err
	}

	// 写入配置文件
	if err := os.WriteFile(cgroupPath, []byte(res.MemoryLimit), 0644); err != nil {
		return err
	}
	return nil
}

func (m MemorySubsystem) Apply(cgroupName string, pid int) error {
	cgroupPath, err := GetCgroupPath(m.Name(), cgroupName)
	if err != nil {
		return err
	}
	log.Infof("%s cgroup path: %s", m.Name(), cgroupName)
	if err != nil {
		return err
	}
	log.Infof("%s cgroup path: %s", m.Name(), cgroupPath)
	// 将 PID 加入当前 Cgroup
	limitFilePath := path.Join(cgroupPath, "tasks")
	if err := os.WriteFile(limitFilePath, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("add pid to cgroup failed: %v", err)
	}
	return nil
}

func (m MemorySubsystem) Remove(cgroupName string) error {
	cgroupPath, err := GetCgroupPath(m.Name(), cgroupName)
	if err != nil {
		return err
	}
	log.Infof("%s cgroup path: %s", m.Name(), cgroupName)
	return os.RemoveAll(cgroupPath)
}
