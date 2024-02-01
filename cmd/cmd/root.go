package cmd

import (
	"log"
	"os"

	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/ui"
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
			err  error
			data []map[string]any
			res  *srch.Response
		)

		vals := FlagsToParams(cmd.Flags())
		idx, err := srch.New(vals)
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case cmd.Flags().Changed(J.Long()):
			j, err := cmd.Flags().GetString(J.Long())
			if err != nil {
				log.Fatal(err)
			}
			data, err = srch.StringSrc(j)
			if err != nil {
				log.Fatal(err)
			}
			idx = idx.Index(data)
		default:
			//in := cmd.InOrStdin()
			//data, err = srch.ReaderSrc(in)
			//if err != nil {
			//log.Fatal(err)
			//}
		}

		res = idx.Search(vals.Encode())

		if cmd.Flags().Changed(B.Long()) {
			tui := ui.New(res.Index)
			res, err = tui.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if res != nil {
				idx = res.Index
			}
		}

		if p, err := cmd.Flags().GetBool("pretty"); err == nil && p {
			idx.PrettyPrint()
		} else {
			println(res.NbHits())
			//idx.Print()
			//d, err := json.Marshal(res)
			//if err != nil {
			//log.Fatal(err)
			//}
			//println(string(d))
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
