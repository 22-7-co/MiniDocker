package runtime

import (
	"encoding/json"
	"fmt"
	"mini-docker/Docker/container"
	"os"
	"strconv"
	"syscall"
)

/*
通过名字停止制定容器
*/
func StopContainer(containerName string) error {
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		return err
	}

	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return fmt.Errorf("convert pid %s to int err: %v", pid, err)
	}

	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		return fmt.Errorf("send sigterm %s, err: %v", pid, err)
	}

	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		return err
	}
	containerInfo.Status = container.STOP
	containerInfo.Pid = ""
	newContarinerInfo, err := json.Marshal(containerInfo)
	if err != nil {
		return fmt.Errorf("json marshal %v,err: %v", containerInfo, err)
	}

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configPath := dirUrl + container.ConfigName
	if err := os.WriteFile(configPath, newContarinerInfo, 0622); err != nil {
		return fmt.Errorf("write file %s err: %v", configPath, err)
	}
	return nil
}
