package cmd

import (
	"log"
	"net/url"
	"os"

	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/ui"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
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
			err      error
			keywords string
			q        = make(url.Values)
			idx      *srch.Index
			data     []map[string]any
		)

		if cmd.Flags().Changed("query") {
			query, err := cmd.Flags().GetString("query")
			if err != nil {
				log.Fatal(err)
			}
			q = srch.NewQuery(query)
			idx = srch.New(q)
			if q.Has("q") {
				keywords = q.Get("q")
			}
		}

		if cmd.Flags().Changed("search") {
			keywords, err = cmd.Flags().GetString("search")
			if err != nil {
				log.Fatal(err)
			}
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

		if cmd.Flags().Changed("filter") {
			filters, err := cmd.Flags().GetStringSlice("filter")
			if err != nil {
				log.Fatal(err)
			}
			filter := srch.NewQuery(lo.ToAnySlice(filters)...)
			for k, vals := range filter {
				for _, v := range vals {
					q.Add(k, v)
				}
			}
		}

		if cmd.Flags().Changed("text") {
			fields, err := cmd.Flags().GetStringSlice("text")
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range fields {
				q.Add("field", f)
			}
		}

		idx = srch.New(q)

		switch {
		case cmd.Flags().Changed("dir"):
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				log.Fatal(err)
			}
			data, err = srch.DirSrc(dir)
			if err != nil {
				log.Fatal(err)
			}
		case len(dataFiles) > 0:
			data, err = srch.FileSrc(dataFiles...)
			if err != nil {
				log.Fatal(err)
			}
		case cmd.Flags().Changed("json"):
			j, err := cmd.Flags().GetString("json")
			if err != nil {
				log.Fatal(err)
			}
			data, err = srch.StringSrc(j)
			if err != nil {
				log.Fatal(err)
			}
		default:
			in := cmd.InOrStdin()
			data, err = srch.ReaderSrc(in)
			if err != nil {
				log.Fatal(err)
			}
		}

		if cmd.Flags().Changed("browse") {
			tui := ui.Browse(q, data)
			idx, err = tui.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			idx = idx.Index(data)

			if keywords != "" {
				idx = idx.Search(keywords)
			}
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

	rootCmd.PersistentFlags().StringSliceVarP(&dataFiles, "index", "i", []string{}, "list of data files to index")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "directory of data files")
	rootCmd.PersistentFlags().StringP("json", "j", "", "json formatted input")

	//rootCmd.MarkFlagsOneRequired("index", "dir", "json")
	rootCmd.MarkFlagsMutuallyExclusive("index", "dir", "json")
	rootCmd.MarkPersistentFlagDirname("dir")
	rootCmd.MarkPersistentFlagFilename("index", ".json")

	rootCmd.PersistentFlags().Bool("ui", false, "select results in a tui")
	rootCmd.PersistentFlags().BoolP("browse", "b", false, "browse results in a tui")

	rootCmd.PersistentFlags().StringSliceP("filter", "f", []string{}, "facet filters")
	rootCmd.PersistentFlags().StringSliceP("text", "t", []string{}, "text fields")
	rootCmd.PersistentFlags().StringP("query", "q", "", "encoded query/filter string (eg. color=red&color=pink&category=post")
	rootCmd.PersistentFlags().StringP("search", "s", "", "search index")
	rootCmd.PersistentFlags().StringSliceP("or", "o", []string{}, "disjunctive facets")
	rootCmd.PersistentFlags().StringSliceP("and", "a", []string{}, "conjunctive facets")

	rootCmd.PersistentFlags().Bool("pretty", false, "pretty print json output")

	rootCmd.PersistentFlags().IntP("workers", "w", 1, "number of workers for computing facets")
	viper.BindPFlag("workers", rootCmd.Flags().Lookup("workers"))

}

func initConfig() {
	viper.AutomaticEnv()
}
