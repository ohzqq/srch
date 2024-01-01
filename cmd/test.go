
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)



var tootCmd = &cobra.Command{
	Use: "toot",
	Aliases: []string{"t"},
	Short: "toot command",
	Long: "long explanation",
	Run: tootCmdRun,
}

func init() {