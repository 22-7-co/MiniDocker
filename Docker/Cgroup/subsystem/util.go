package subsystem

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

// 获取Cgroup的路径
func GetCgroupPath(subsystem string, cgroupName string) (string, error) {
	cgroupRoot, err := FindCgroupMountPoint(subsystem)
	if err != nil {
		return "", err
	}
	cgroupPath := path.Join(cgroupRoot, cgroupName)
	_, err = os.Stat(cgroupPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("file stat err: %v", err)
	}
	if os.IsNotExist(err) {
		if err := os.Mkdir(cgroupPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("mkdir err: %v", err)
		}
	}
	return cgroupPath, nil
}

/*
在/proc/self/mountinfo中查找某个 cgroup 子系统的挂载信息
*/
func FindCgroupMountPoint(subsystem string) (string, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", fmt.Errorf("open /proc/self/mountinfo err: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.SplitN(txt, " ", 5)
		if len(fields) < 5 {
			continue
		}

		mountPoint := fields[4]
		options := fields[len(fields)-1] // 最后是选项

		// cgroup v2 特征：类型是 cgroup2，或者选项包含 cgroup2
		if strings.Contains(options, "cgroup2") || strings.HasPrefix(mountPoint, "/sys/fs/cgroup") {
			// v2 统一挂载点，直接返回（不管 subsystem 是什么）
			return "/sys/fs/cgroup", nil
		}

		// v1 原来的逻辑
		log.Debugf("mount info txt fields: %s", fields)
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4], nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("file scanner err: %v", err)
	}
	return "", fmt.Errorf("FindCgroupMountPoint is empty")
}
