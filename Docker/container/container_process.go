package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func NewParentProcess(tty bool, containerName, rootUrl, mntUrl, volume string, envSlice []string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorf("create pipe err: %v", err)
		return nil, nil
	}
	/*
		创建一个新进程，执行当前程序自身（/proc/self/exe 是当前运行程序的路径）
		相当于：让当前程序（mydocker）再启动一个子进程，执行 "mydocker init command"
		这是容器实现的经典技巧：用同一个二进制文件，通过不同命令（run/init）实现不同逻辑
	*/
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // UTS隔离：容器有独立的主机名、域名
			syscall.CLONE_NEWPID | // PID隔离：容器内PID从1开始（比如容器内ps aux只有自己的进程）
			syscall.CLONE_NEWNS | // 挂载点隔离：容器内的文件系统挂载不影响宿主机
			syscall.CLONE_NEWNET | // 网络隔离：容器有独立的网络栈（网卡、IP、端口）
			syscall.CLONE_NEWIPC, // IPC隔离：容器内的IPC（信号量、消息队列）不与宿主机互通
	}
	/*
		如果开启tty（交互模式，对应 docker run -ti），
		则把父进程的标准输入/输出/错误重定向给子进程
	*/
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		dirUrl := fmt.Sprintf(DefaultInfoLocation, containerName)
		if err := os.MkdirAll(dirUrl, 0622); err != nil {
			log.Errorf("mkdir %s err: %v", dirUrl, err)
			return nil, nil
		}

		stdLogFilePath := dirUrl + ContainerLogFile
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorf("create file %s err: %v", stdLogFilePath, err)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}

	// 父进程和子进程之间的通信
	cmd.ExtraFiles = []*os.File{readPipe}
	if err := newWorkSpace(rootUrl, mntUrl, volume, containerName); err != nil {
		log.Errorf("new work space err: %v", err)
		return nil, nil
	}
	cmd.Dir = mntUrl
	cmd.Env = append(os.Environ(), envSlice...)
	return cmd, writePipe
}

func newWorkSpace(rootUrl, mntUrl, volume, containerName string) error {
	if err := createReadOnlyLayer(rootUrl); err != nil {
		return err
	}
	if err := createWriteLayer(containerName); err != nil {
		return err
	}
	if err := createMountPoint(rootUrl, mntUrl, containerName); err != nil {
		return err
	}
	if err := mountExtractVolume(mntUrl, volume, containerName); err != nil {
		return err
	}
	return nil
}

// 创建一个名为busybox 的文件夹作为容器唯一只读层
func createReadOnlyLayer(busyboxUrl string) error {
	exist, err := pathExist(busyboxUrl)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("busy dir don't exist: %s", busyboxUrl)
	}
	return nil
}

func mountExtractVolume(mntUrl, volume, containerName string) error {
	if volume == "" {
		return nil
	}
	volumeUrls := strings.Split(volume, ":")
	length := len(volumeUrls)
	if length != 2 || volumeUrls[0] == "" || volumeUrls[1] == "" {
		return fmt.Errorf("volume parameter input is invalid")
	}
	return mountVolume(mntUrl+containerName+"/", volumeUrls)
}

func mountVolume(mntUrl string, volumeUrls []string) error {
	// 如果目录不存在，创建目录
	parentUrl := volumeUrls[0]
	exist, err := pathExist(parentUrl)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if !exist {
		// MkdirAll 递归创建多级目录
		if err := os.MkdirAll(parentUrl, 0777); err != nil {
			return fmt.Errorf("create parent dir failed: %v", err)
		}
	}

	// 容器内创建挂载点
	containerUrl := mntUrl + volumeUrls[1]
	if err := os.MkdirAll(containerUrl, 0777); err != nil {
		return fmt.Errorf("mkdir container volume err: %v", err)
	}

	// 把宿主机文件目录挂载到容器
	dirs := "dirs=" + parentUrl
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mount volume err: %v", err)
	}
	return nil
}

// 创建一个名为writeLayer 的文件夹作为容器唯一可写层
func createWriteLayer(containerName string) error {
	writeUrl := RootUrl + "/witeLayer/" + containerName + "/"
	exist, err := pathExist(writeUrl)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if !exist {
		if err := os.MkdirAll(writeUrl, 0777); err != nil {
			return fmt.Errorf("create write layer failed: %v", err)
		}
	}
	return nil
}

func createMountPoint(rootUrl, mntUrl, containerName string) error {
	// 创建一个名为mnt 的文件夹作为容器唯一挂载点
	mountPath := mntUrl + containerName + "/"
	log.Infof("root url: %s, mutUrl: $s, mountPath: %s", rootUrl, mntUrl, mountPath)
	exist, err := pathExist(mountPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if !exist {
		if err := os.MkdirAll(mountPath, 0777); err != nil {
			return fmt.Errorf("create mount point failed: %v", err)
		}
	}
	// 将 writeLayer, busybox 挂载到 mnt
	dirs := "dirs=" + RootUrl + "/writeLayer/" + rootUrl
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mountPath)
	log.Infof(cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mnt dir err: %v", err)
	}
	return nil
}

func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return false, err
}
