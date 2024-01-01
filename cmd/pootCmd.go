package main

import (
	"github.com/spf13/cobra"
)

func pootCmdRun(cmd *cobra.Command, args []string) {
	println(cmd.Name())
}
