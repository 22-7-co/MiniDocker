package runtime

import (
	"encoding/json"
	"fmt"
	"mini-docker/Docker/container"
	"os"
	"strings"
)

func getContainerPidByName(containerName string) (string, error) {
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirUrl + container.ConfigName

	contentBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return "", fmt.Errorf("read file %s err: %v", configFilePath, err)
	}

	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", fmt.Errorf("json unmarshal err: %v", err)
	}
	return containerInfo.Pid, nil
}

func getContainerInfoByName(containerName string) (*container.ContainerInfo, error) {
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirUrl + container.ConfigName
	contentBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("read file %s err: %v", configFilePath, err)
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}
	return &containerInfo, nil
}

func getEnvsByPid(pid string) ([]string, error) {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s err: %v", path, err)
	}
	return strings.Split(string(contentBytes), "\u0000"), nil
}
