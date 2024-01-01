package main

import (
	"github.com/spf13/cobra"
)

func tootCmdRun(cmd *cobra.Command, args []string) {
	println(cmd.Name())
}
