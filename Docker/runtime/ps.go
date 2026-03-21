package runtime

import (
	"encoding/json"
	"fmt"
	"mini-docker/Docker/container"
	"os"
)

func ListContainers() error {
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, "")
	dirUrl = dirUrl[:len(dirUrl)-1]
	files, err := os.ReadDir(dirUrl)
	if err != nil {
		return fmt.Errorf("read dir %s err: %v", dirUrl, err)
	}

	var containers []*container.ContainerInfo
	for _, file := range files {
		tmpContainer, err := getContainerInfo(file)
		if err != nil {
			return err
		}
		containers = append(containers, tmpContainer)
	}

	// w :=
	return nil
}

func getContainerInfo(file os.DirEntry) (*container.ContainerInfo, error) {
	containerName := file.Name()
	configFileDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := configFileDir + container.ConfigName

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("read file %s err %v", configFilePath, err)
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}
	return &containerInfo, nil
}
