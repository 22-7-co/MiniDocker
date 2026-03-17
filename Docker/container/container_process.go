package container

import (
	"fmt"
	"os"
	"os/exec"
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
	// if err := createReadOnlyLayer(rootUrl); err != nil {
	// 	return err
	// }
	// if err := createWriteLayer(containerName); err != nil {
	// 	return err
	// }
	// if err := createMountPoint(rootUrl, mntUrl, containerName); err != nil {
	// 	return err
	// }
	// if err := mountExtractVolume(mntUrl, volume, containerName); err != nil {
	// 	return err
	// }
	return nil
}
