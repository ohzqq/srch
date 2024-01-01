package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tootCmd = &cobra.Command{
	Use:     "toot",
	Aliases: []string{"t"},
	Short:   "toot command",
	Long:    "long explanation",
	Run:     tootCmdRun,
}

func init() {
	rootCmd.AddCommand(tootCmd)

	tootCmd.Flags().StringP(
		"turkey",
		"t",
		"",
		"gobble gobble",
	)

	viper.BindPFlag("turkey", tootCmd.Flags().Lookup("turkey"))

	tootCmd.PersistentFlags().StringP(
		"gooble",
		"g",
		"",
		"gobble gobble",
	)

}

var pootCmd = &cobra.Command{
	Use: "poot",

	Run: pootCmdRun,
}

func init() {

}
