package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ohzqq/srch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	dataFiles []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "srch -f file... | -d dir | -i string [flags]",
	Short: "search collections",
	Long: `facet aggregates data on specified fields, with option filters. 

The command accepts stdin, flags, and positional arguments.

If a config file has a "data" field no other argument or flag is required. 

Without the "data" field, data must be specified through a flag or positional
argument.

By default, results are printed to stdout as json.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.Lshortfile)

		var (
			err     error
			filters string
			q       = make(url.Values)
			idx     *Index
		)

		if cmd.Flags().Changed("ui") {
			opts = append(opts, srch.Interactive)
		}

		if cmd.Flags().Changed("query") {
			query, err = cmd.Flags().GetString("query")
			if err != nil {
				log.Fatal(err)
			}
			q = srch.NewQuery(query)
		}

		if cmd.Flags().Changed("or") {
			or, err := cmd.Flags().GetStringSlice("or")
			if err != nil {
				log.Fatal(err)
			}
			for _, o := range or {
				q.Add("or", o)
			}
		}

		if cmd.Flags().Changed("and") {
			and, err := cmd.Flags().GetStringSlice("and")
			if err != nil {
				log.Fatal(err)
			}
			for _, o := range and {
				q.Add("and", o)
			}
		}

		idx = srch.New(q)

		switch {
		case cmd.Flags().Changed("dir"):
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				log.Fatal(err)
			}
			idx.Index(srch.DirSrc(dir))
		case len(dataFiles) > 0:
			idx.Index(srch.FileSrc(dataFiles...))
		case cmd.Flags().Changed("json"):
			j, err := cmd.Flags().GetString("json")
			if err != nil {
				log.Fatal(err)
			}
			idx.Index(srch.StringSrc(j))
		default:
			in := cmd.InOrStdin()
			idx.Index(srch.ReaderSrc(in))
		}

		if cmd.Flags().Changed("search") {
			kw, err := cmd.Flags().GetString("search")
			if err != nil {
				log.Fatal(err)
			}
			idx = idx.Search(kw)
		}

		if p, err := cmd.Flags().GetBool("pretty"); err == nil && p {
			idx.PrettyPrint()
		} else {
			idx.Print()
		}
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

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "json formatted config file")

	rootCmd.PersistentFlags().StringSliceVarP(&dataFiles, "index", "i", []string{}, "list of data files to index")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "directory of data files")
	rootCmd.PersistentFlags().StringP("json", "j", "", "json formatted input")

	rootCmd.MarkFlagsOneRequired("index", "dir", "json")
	rootCmd.MarkFlagsMutuallyExclusive("index", "dir", "json")
	rootCmd.MarkPersistentFlagDirname("dir")
	rootCmd.MarkPersistentFlagFilename("index", ".json")

	rootCmd.PersistentFlags().Bool("ui", false, "select results in a tui")

	rootCmd.PersistentFlags().StringP("query", "q", "", "encoded query/filter string (eg. color=red&color=pink&category=post")
	rootCmd.PersistentFlags().StringP("search", "s", "", "search index")
	rootCmd.PersistentFlags().StringSliceP("or", "o", "", "disjunctive facets")
	rootCmd.PersistentFlags().StringSliceP("and", "a", "", "conjunctive facets")

	rootCmd.PersistentFlags().Bool("pretty", false, "pretty print json output")

	rootCmd.PersistentFlags().IntP("workers", "w", 1, "number of workers for computing facets")
	viper.BindPFlag("workers", rootCmd.Flags().Lookup("workers"))

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
