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

	// 由于没有隔离文件系统，所以需要在容器中重新挂载一下/proc，否则ps等命令还是会读取到根命名空间中的进程
	fmt.Println("挂载proc目录")
	syscall.Mount("proc", "/proc", "proc", 0, "")

	// 使用此系统调用来执行命令，执行的命令将会取代父进程(称为隔离空间中的1号进程)
	fmt.Println("执行命令:", args[0])
	syscall.Exec(args[0], args[0:], os.Environ())

	syscall.Unmount("/proc", 0)
}
