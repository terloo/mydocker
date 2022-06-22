package cmd

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
)

var childCmd = &cobra.Command{
	Use: "child",
	Run: childFunc,
}

// mydocker child [newHostname]
func childFunc(cmd *cobra.Command, args []string) {
	// 修改hostname
	fmt.Println("修改hostname为", args[0])
	syscall.Sethostname([]byte(args[0]))
}
