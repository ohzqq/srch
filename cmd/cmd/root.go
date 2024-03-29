package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ohzqq/srch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "srch -j json... | -d dir | -b string [flags]",
	Short: "search collections",
	Long: `srch searches a collection.

The command accepts stdin, flags, and positional arguments.

If a config file has a "data" field no other argument or flag is required. 

Without the "data" field, data must be specified through a flag or positional
argument.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Lshortfile)

		req := srch.GetViperParams()

		idx, err := srch.New(req.String())
		if err != nil {
			println(req.String())
			log.Fatal(err)
		}

		res, err := idx.Search(req.String())
		if err != nil {
			log.Fatal(err)
		}
		println(res.NbHits)

		d, err := json.Marshal(res)
		if err != nil {
			log.Fatal(err)
		}

		println(string(d))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	defineFlags()

	//rootCmd.MarkFlagsOneRequired("index", "dir", "json")
	//rootCmd.MarkFlagsMutuallyExclusive("blv", "dir", "json")
	rootCmd.MarkPersistentFlagDirname("dir")
	rootCmd.MarkPersistentFlagFilename("json", ".json", ".ndjson")

	rootCmd.PersistentFlags().
		Bool(
			"ui",
			false,
			"select results in a tui",
		)

	rootCmd.PersistentFlags().
		Bool(
			"pretty",
			false,
			"pretty print json output",
		)

	rootCmd.PersistentFlags().
		IntP(
			"workers",
			"w",
			4,
			"number of workers for computing facets",
		)
	viper.BindPFlag(
		"workers",
		rootCmd.Flags().Lookup("workers"),
	)

}

func initConfig() {
	viper.AutomaticEnv()

}
