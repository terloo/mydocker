package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "mydocker",
}

func init() {
	RootCmd.AddCommand(runCmd)
	RootCmd.AddCommand(childCmd)
}