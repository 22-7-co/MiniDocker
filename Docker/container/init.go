package container

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

/*
RunContainerInitProcess 运行在容器内部的第一个进程（PID 1） 中：
先挂载 proc 文件系统（容器内必须挂载，否则无法使用 ps/top 等查看进程的命令）；
执行用户传入的命令（比如 echo hello、sh 等），这也是容器最终要运行的核心逻辑。
*/

func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s, args %s", command, args)
	defaultMountFlags := syscall.MS_NOEXEC | //禁止在 proc 分区执行可执行文件（安全）
		syscall.MS_NOSUID | //MS_NOSUID：禁止 setuid/setgid 权限（防止提权）
		syscall.MS_NODEV //MS_NODEV：禁止访问设备文件（安全）

	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		return err
	}

	argv := []string{command}
	// syscall.Exec 会用新的命令替换当前进程的镜像，当前进程的代码会被替换，不再执行后续逻辑
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}
