package runtime

import (
	"encoding/json"
	"fmt"
	"math/rand"
	cgroup "mini-docker/Docker/Cgroup"
	"mini-docker/Docker/Cgroup/subsystem"
	"mini-docker/Docker/container"
	"mini-docker/Docker/network"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Run
/*
这里的Start方法是真正开始前面创建好的 command 的调用，
它首先会clone出来一个namespace隔离的进程，然后在子进程中，调用/proc/self/exe,也就是自己调用自己
发送 init 参数，调用我们写的 init 方法，去初始化容器的一些资源
*/
func Run(tty, detach bool, cmdArray []string, config *subsystem.ResourceConfig, volume, containerName string,
	envSlice []string, nw string, portMapping []string) {
	id, containerName := getContainerName(containerName)

	mntUrl := container.RootUrl + "/mnt/"
	rootUrl := container.BusyboxPath
	parent, writePipe := container.NewParentProcess(tty, containerName, rootUrl, mntUrl, volume, envSlice)

	if err := parent.Start(); err != nil {
		log.Errorf("start parent process err: %v", err)

		deleteWorkSpace(rootUrl, mntUrl, volume, containerName)
		return
	}

	// 记录容器信息
	containerName, err := recordContainerInfo(parent.Process.Pid, cmdArray, id, containerName)
	if err != nil {
		log.Error(err)
		return
	}

	cgroupMgr := cgroup.NewCgroupManager("mydocker-cgroup")
	defer cgroupMgr.Destroy()
	if err != nil {
		log.Errorf("cgroup apply err: %v", err)
		return
	}

	if nw != "" {
		// 定义网络
		_ = network.Init()
		containerInfo := &container.ContainerInfo{
			ID:           id,
			Pid:          strconv.Itoa(parent.Process.Pid),
			Name:         containerName,
			PortMappings: portMapping,
		}
		if err := network.Connect(nw, containerInfo); err != nil {
			log.Errorf("connect network err: %v", err)
			return
		}
	}

	sendInitCommand(cmdArray, writePipe)

	log.Infof("parent  process run")
	if !detach {
		_ = parent.Wait()
		deleteWorkSpace(rootUrl, mntUrl, volume, containerName)
		deleteContainerInfo(containerName)
	}
	os.Exit(-1)
}

func deleteContainerInfo(containerName string) {
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirUrl); err != nil {
		log.Errorf("remove container info dir err: %v", err)
	}
}

// 为了防止命名冲突，为 容器生成一个随机的名字
func getContainerName(containerName string) (string, string) {
	id := randStringBytes(10)
	if containerName == "" {
		containerName = id
	}
	return id, containerName
}

func recordContainerInfo(pid int, cmdArray []string, id, containerName string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(cmdArray, " ")
	containerInfo := container.ContainerInfo{
		ID:         id,
		Pid:        strconv.Itoa(pid),
		Name:       containerName,
		Command:    command,
		CreateTime: createTime,
		Status:     container.RUNNING,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		return "", fmt.Errorf("container info to json string err: %v", err)
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		return "", fmt.Errorf("mkdir %s err : %v", dirUrl, err)
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return "", fmt.Errorf("create file %s, err: %v", fileName, err)
	}

	if _, err := file.WriteString(jsonStr); err != nil {
		return "", fmt.Errorf("file write string err: %v", err)
	}
	return containerName, nil
}

// 生成一个指长度的纯数字随机字符串
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func sendInitCommand(arrary []string, writePipe *os.File) {
	command := strings.Join(arrary, " ")
	log.Infof("all command is : %s", command)

	if _, err := writePipe.WriteString(command); err != nil {
		log.Errorf("write pipe write string err : %v", err)
		return
	}

	if err := writePipe.Close(); err != nil {
		log.Errorf("write pipe close err : %v", err)
	}
}

func deleteWorkSpace(rootUrl, mntUrl, volume, containerName string) {
	umountVolume(mntUrl, volume, containerName)
	deleteMountPoint(mntUrl + containerName + "/")
	deleteWriteLayer(rootUrl, containerName)
}

func umountVolume(mntUrl, volume, containerName string) {
	if volume == "" {
		return
	}
	volumeUrls := strings.Split(volume, ":")
	if len(volumeUrls) != 2 || volumeUrls[0] == "" || volumeUrls[1] == "" {
		return
	}

	containerUrl := mntUrl + containerName + "/" + volumeUrls[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("umount volume err: %v", err)
	}
}

func deleteMountPoint(mntUrl string) {
	cmd := exec.Command("umount", mntUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("deleteMountPoint umount %s err : %v", mntUrl, err)
	}
	if err := os.RemoveAll(mntUrl); err != nil {
		log.Errorf("delete mnt dir %s err: %v", mntUrl, err)
	}
}

func deleteWriteLayer(rootUrl, containerName string) {
	wirteUrl := rootUrl + "writeLayer/" + containerName
	if err := os.RemoveAll(wirteUrl); err != nil {
		log.Errorf("delete write layer err: %v", err)
	}
}
