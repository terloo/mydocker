package cmd

import (
	"fmt"
	"os"
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

	fmt.Println("修改hostname为mycontainer")
	syscall.Sethostname([]byte("mycontainer"))

	// 改变当前进程所使用的根目录，此目录需要把容器执行命令的所需的文件拷贝进去
	syscall.Chroot("rootfs")
	syscall.Chdir("/")

	// 挂载/proc，此时/proc会被挂载到rootfs中
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		// 如果没有/proc文件夹，会panic
		panic(err)
	}

	// 使用此系统调用来执行命令，执行的命令将会取代父进程(称为隔离空间中的1号进程)
	fmt.Println("执行命令:", args[0])
	if err := syscall.Exec(args[0], args[0:], os.Environ()); err != nil {
		// 如果没有命令文件，会panic
		panic(err)
	}

	syscall.Unmount("/proc", 0)
}
