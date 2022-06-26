package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use: "init",
	Run: initFunc,
}

// mydocker init /bin/bash [其他参数]
// 此命令已在隔离环境中进行执行
func initFunc(cmd *cobra.Command, args []string) {

	// 将镜像和容器目录放在指定位置，每次启动容器时，将镜像里的文件拷贝到容器对应目录中
	imageFolderPath := "./base"
	rootFolderPath := "./rootfs"

	// 如果容器目录不存在，则将镜像里的文件拷贝到容器目录
	if _, err := os.Stat(rootFolderPath); os.IsNotExist(err) {
		if err := CopyFileOrDirectory(imageFolderPath, rootFolderPath); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	if err := syscall.Sethostname([]byte("mycontainer")); err != nil {
		panic(err)
	}

	// 改变当前进程所使用的根目录，此目录需要把容器执行命令的所需的文件拷贝进去
	if err := syscall.Chroot("rootfs"); err != nil {
		panic(err)
	}
	if err := syscall.Chdir("/"); err != nil {
		panic(err)
	}

	// 挂载/proc，此时/proc会被挂载到rootfs中
	if err := os.Mkdir("/proc", os.ModeDir); !os.IsExist(err) && err != nil {
		panic(err)
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		// 如果没有/proc文件夹，会panic
		panic(err)
	}

	// 在隔离环境下搜索可执行文件的路径
	cmdPath, err := exec.LookPath(args[0])
	fmt.Println("执行命令:", cmdPath)
	if err != nil {
		panic(err)
	}

	// 使用此系统调用来执行命令，执行的命令将会取代父进程(称为隔离空间中的1号进程)
	// 注意第二个参数将会被保存于/proc/<pid>/cmdline，来表示此命名执行的信息(命令，参数，选项等)
	if err := syscall.Exec(cmdPath, args, os.Environ()); err != nil {
		// 如果没有命令文件，会panic
		panic(err)
	}

	// TODO Exec后面的代码无法执行
	fmt.Println("取消挂载")
	// 命令执行结束后取消/proc的挂载
	if err := syscall.Unmount("/proc", 0); err != nil {
		panic(err)
	}

}

func CopyFileOrDirectory(srcFile, destFile string) error {
	cpCmd := exec.Command("cp", "-r", srcFile, destFile)
	return cpCmd.Run()
}
