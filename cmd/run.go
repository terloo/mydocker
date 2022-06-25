package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use: "run",
	Run: runFunc,
}

// mydocker run /bin/bash [其他参数]
// run命令使用系统调用调用自身的init命令，来创建一个隔离进程
func runFunc(cmd *cobra.Command, args []string) {
	// os.Args[0]为本命令，即mydocker
	// init 为子命令
	// args[1:]为其他参数
	// 最终执行的命令为mydocker init /bin/bash [其他参数]
	initArgs := append([]string{"init"}, args...)
	initCmd := exec.Command(os.Args[0], initArgs...)
	// 使用linux提供的系统调用来修改该命令执行时的隔离级别
	// 这里隔离UTS(hostname domainname)，PID
	initCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	initCmd.Stdin = os.Stdin
	initCmd.Stderr = os.Stderr
	initCmd.Stdout = os.Stdout
	if err := initCmd.Run(); err != nil {
		panic(err)
	}
}
