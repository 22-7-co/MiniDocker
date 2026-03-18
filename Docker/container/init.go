package container

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// 读取用户传入的参数
func readUserCommand() []string {
	readPipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(readPipe)
	if err != nil {
		log.Errorf("read int argv pipe err: %v", err)
		return nil
	}
	return strings.Split(string(msg), " ")
}

// 初始挂载点
func setUpMount() error {
	// 首先设置根目录为私有模式，防止影响pivot_root
	if err := syscall.Mount("/", "/", "", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		return fmt.Errorf("setUpMount mount proc err: %v", err)
	}

	// 进入Busybox，固定路径，busybox提前解压好，放到指定目录
	err := privotRoot(BusyboxPath)
	if err != nil {
		return err
	}

	// mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		log.Errorf("proc mount failed: %v", err)
		return err
	}
	syscall.Mount("tmpfs", "/dev", "tempfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	return nil
}

func privotRoot(root string) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("pwd err: %v", err)
		return err
	}
	log.Infof("current pwd: %s", pwd)

	// 为了使当前root的老root和新root不在同一个文件系统下，我们把root重新mount一次
	// bind mount 是把相同的内容换了一个挂载点的挂载方法
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}
	// 创建 rootfs, .pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if _, err := os.Stat(pivotDir); err == nil {
		if err := os.Remove(pivotDir); err != nil {
			return err
		}
	}

	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir of pivot_root err:%v", err)
	}

	// pivot_root 到新的rootfs, 老的old_root 现在挂载在rootfs/.pivot_root 上
	// 挂载点目前依然可以在mount中看到
	log.Infof("root: %s， pivotDir: %s", root, pivotDir)
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root err: %v", err)
	}
	// 修改工作目录
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir err: %v", err)
	}

	// 取消临时文件 .pivot_root的挂载并删除它
	// 注意当前已经在根目录下，所有临时文件的目录也改变了
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("umount pivot_root err: %v", err)
	}
	return os.Remove(pivotDir)
}
