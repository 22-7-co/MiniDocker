package runtime

import (
	"fmt"
	"mini-docker/Docker/container"
	"os"
)

/*
通过指定容器名移除容器
*/
func RemoveContainer(containerName string) error {
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		return err
	}

	if containerInfo.Status != container.STOP {
		return fmt.Errorf("could remove running container")
	}

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirUrl); err != nil {
		return fmt.Errorf("remove dir %s, err: %v", dirUrl, err)
	}
	return nil
}
