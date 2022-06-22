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

// mydocker run bash
func runFunc(cmd *cobra.Command, args []string) {
	// 创建一个操作系统命令，此demo传入的命令为bash，方便观察
	systemCmd := exec.Command(args[0])
	// 使用linux提供的系统调用来修改该命令执行时的隔离级别
	// 这里隔离UTS(hostname domainname)
	systemCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	systemCmd.Stdin = os.Stdin
	systemCmd.Stderr = os.Stderr
	systemCmd.Stdout = os.Stdout
	if err := systemCmd.Run(); err != nil {
		panic(err)
	}
}
