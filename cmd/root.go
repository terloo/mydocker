package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "mydocker [Command]",
}

func init() {
	RootCmd.AddCommand(runCmd)
	RootCmd.AddCommand(initCmd)
	// RootCmd.AddCommand(logsCmd)
}