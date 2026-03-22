package runtime

import (
	"fmt"
	"mini-docker/Docker/container"
	"os"
)

func LogContainer(containerName string) error {
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	logFilePath := dirUrl + container.ContainerLogFile
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("open file %s err: %v", logFilePath, err)
	}
	defer file.Close()

	content, err := os.ReadFile(logFilePath)
	if err != nil {
		return fmt.Errorf("read file %s err: %v", logFilePath, err)
	}
	fmt.Println(string(content))
	return nil
}
