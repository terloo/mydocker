package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command in a new container",
	Run:   runFunc,
}

func init() {
	// -i 是否为交互式运行
	runCmd.Flags().BoolP("interactive", "i", false, "Keep STDIN open even if not attached")
	// -t 是否模拟一个tty
	runCmd.Flags().BoolP("tty", "t", false, "Allocate a pseudo-TTY")
	// -d 是否在后台运行
	runCmd.Flags().BoolP("detach", "d", false, "Run container in background and print container ID")
}

// mydocker run /bin/bash [其他参数]
// run命令使用系统调用调用自身的init命令，来创建一个隔离进程
func runFunc(cmd *cobra.Command, args []string) {
	isTTY, _ := cmd.Flags().GetBool("tty")
	isInteractive, _ := cmd.Flags().GetBool("interactive")

	// /proc/self/exe为本命令，即mydocker
	// init 为子命令
	// args[1:]为其他参数
	// 最终执行的命令为mydocker init /bin/bash [其他参数]
	initArgs := append([]string{"init"}, args...)
	initCmd := exec.Command("/proc/self/exe", initArgs...)

	// 使用linux提供的系统调用来修改init子命令执行时的隔离级别
	initCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
		// 隔离USER时，需要指定映射配置  容器外的id->容器内的id
		UidMappings: []syscall.SysProcIDMap{
			{
				// 容器内的用户id，0为root
				ContainerID: 0,
				// 容器外的用户id，获取当前用户id
				HostID: os.Getuid(),
				Size:   1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				// 容器内使用的gid，0为root组
				ContainerID: 0,
				// 容器外的组id，获取当前组id
				HostID: os.Getgid(),
				Size:   1,
			},
		},
	}
	// 模拟一个tty
	if isTTY {
		// 交互式，可以输入
		if isInteractive {
			initCmd.Stdin = os.Stdin
		}
		initCmd.Stderr = os.Stderr
		initCmd.Stdout = os.Stdout
	}
	if err := initCmd.Start(); err != nil {
		panic(err)
	}
	initCmd.Wait()
}
