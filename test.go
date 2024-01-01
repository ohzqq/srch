
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}



var rootCmd = &cobra.Command{
	Use: "test",
	
	Short: "test cmd",
	
	Run: rootCmdRun,
}

func init() {
		rootCmd.AddCommand(tootCmd)
		rootCmd.AddCommand(pootCmd)

	
	rootCmd.PersistentFlags().String(
		"config",
		
		"",
		"",
	)
		
		
}
		

var tootCmd = &cobra.Command{
	Use: "toot",
	Aliases: []string{"t"},
	Short: "toot command",
	Long: "long explanation",
	Run: tootCmdRun,
}

func init() {

	
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
		

