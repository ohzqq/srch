package main

import (
	"github.com/spf13/cobra"
)

func rootCmdRun(cmd *cobra.Command, args []string) {
	println(cmd.Name())
}
