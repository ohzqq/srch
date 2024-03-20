package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/param"
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

By default, results are printed to stdout as json.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Lshortfile)

		//b, err := cmd.Flags().GetStringSlice("blv")
		//if err != nil {
		//log.Fatal(err)
		//}
		//fmt.Printf("viper %#v\n", b)

		req := srch.GetViperParams()
		println(req.String())
		idx, err := srch.New(req.String())
		if err != nil {
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
		//var (
		//  err  error
		//  data []map[string]any
		//  res  *srch.Response
		//)

		//if p, err := cmd.Flags().GetBool("pretty"); err == nil && p {
		//  res.PrettyPrint()
		//} else {
		//  res.Print()
		//}

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
	for _, key := range param.SettingParams {
		switch key {
		case param.SrchAttr:
			viper.SetDefault(key, []string{"title"})
		case param.FacetAttr:
			viper.SetDefault(key, []string{"tags"})
		case param.SortAttr:
			viper.SetDefault(key, []string{"title:desc"})
		case param.UID:
			viper.SetDefault(key, "id")
		}
	}

	for _, key := range param.SearchParams {
		switch key {
		case param.SortFacetsBy:
			viper.SetDefault(key, "tags:count:desc")
		case param.Facets:
			viper.SetDefault(key, []string{"tags"})
		case param.RtrvAttr:
			viper.SetDefault(key, "*")
		case param.Page:
			viper.SetDefault(key, 0)
		case param.HitsPerPage:
			viper.SetDefault(key, -1)
		case param.SortBy:
			viper.SetDefault(key, "title")
		case param.Order:
			viper.SetDefault(key, "desc")
		}
	}

	cobra.OnInitialize(initConfig)

	defineFlags()

	//rootCmd.MarkFlagsOneRequired("index", "dir", "json")
	//rootCmd.MarkFlagsMutuallyExclusive("blv", "dir", "json")
	rootCmd.MarkPersistentFlagDirname("dir")
	rootCmd.MarkPersistentFlagFilename("json", ".json")

	rootCmd.PersistentFlags().
		Bool(
			UI.Long(),
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
			Workers.Long(),
			Workers.Short(),
			1,
			"number of workers for computing facets",
		)
	viper.BindPFlag(
		Workers.Long(),
		rootCmd.Flags().Lookup("workers"),
	)

}

func initConfig() {
	viper.AutomaticEnv()

}
