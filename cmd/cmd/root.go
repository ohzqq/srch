package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ohzqq/facet"
	"github.com/ohzqq/srch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	dataFiles []string
	idx       = &facet.Index{}
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
			opts    []srch.Opt
			src     srch.Src
			q       = make(srch.Query)
		)

		switch {
		case cmd.Flags().Changed("dir"):
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				log.Fatal(err)
			}
			m, err := filepath.Glob(filepath.Join(dir, "/*"))
			if err != nil {
				log.Fatal(err)
			}
			src = srch.FileSrc(m...)
		case len(dataFiles) > 0:
			src = srch.FileSrc(dataFiles...)
		case cmd.Flags().Changed("json"):
			j, err := cmd.Flags().GetString("json")
			if err != nil {
				log.Fatal(err)
			}
			src = srch.ReadDataSrc(bytes.NewBufferString(j))
		default:
			in := cmd.InOrStdin()
			src = srch.ReadDataSrc(in)
		}

		if cfgFile != "" {
			opts = append(opts, srch.CfgFile(cfgFile))
		}

		if cmd.Flags().Changed("ui") {
			opts = append(opts, srch.Interactive)
		}

		idx := srch.New(src, opts...)

		if cmd.Flags().Changed("filter") {
			filters, err = cmd.Flags().GetString("filter")
			if err != nil {
				log.Fatal(err)
			}
			q, err = srch.NewQuery(filters)
			if err != nil {
				log.Fatal(err)
			}
		}

		if cmd.Flags().Changed("search") {
			kw, err := cmd.Flags().GetString("search")
			if err != nil {
				log.Fatal(err)
			}
			q.Set("q", kw)
		}

		if q.String() != "" {
			idx = idx.Search(q.Values())
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

	rootCmd.PersistentFlags().StringP("filter", "f", "", "encoded query/filter string (eg. color=red&color=pink&category=post")
	rootCmd.PersistentFlags().StringP("search", "s", "", "search index")

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
