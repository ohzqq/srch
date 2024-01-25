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

		if cmd.Flags().Changed(P.Long()) {
			query, err := cmd.Flags().GetString(P.Long())
			if err != nil {
				log.Fatal(err)
			}
			q = srch.ParseQuery(query)
			idx, err = srch.New(q)
			if err != nil {
				log.Fatal(err)
			}
			if q.Has(srch.Query) {
				keywords = q.Get(srch.Query)
			}
		}

		if cmd.Flags().Changed(Q.Long()) {
			keywords, err = cmd.Flags().GetString(Q.Long())
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

		if cmd.Flags().Changed(R.Long()) {
			filters, err := cmd.Flags().GetStringSlice(R.Long())
			if err != nil {
				log.Fatal(err)
			}
			filter := srch.ParseQuery(lo.ToAnySlice(filters)...)
			for k, vals := range filter {
				for _, v := range vals {
					q.Add(k, v)
				}
			}
		}

		if cmd.Flags().Changed(S.Long()) {
			fields, err := cmd.Flags().GetStringSlice(S.Long())
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range fields {
				q.Add("field", f)
			}
		}

		idx, err = srch.New(q)
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case cmd.Flags().Changed(D.Long()):
			dir, err := cmd.Flags().GetString(D.Long())
			if err != nil {
				log.Fatal(err)
			}
			data, err = srch.DirSrc(dir)
			if err != nil {
				log.Fatal(err)
			}
		case cmd.Flags().Changed(I.Long()):
			files, err := cmd.Flags().GetStringSlice(I.Long())
			if err != nil {
				log.Fatal(err)
			}
			data, err = srch.FileSrc(files...)
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

		var res *srch.Response
		if cmd.Flags().Changed("browse") {
			tui := ui.Browse(q, data)
			idx, err = tui.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			idx = res.Index.Index(data)

			if keywords != "" {
				res = idx.Search(keywords)
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
}

func initConfig() {
	viper.AutomaticEnv()
}
