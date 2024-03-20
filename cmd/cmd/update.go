package cmd

import (
	"log"
	"path/filepath"

	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/data"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update INDEX",
	Short: "update documents in an index",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		path := args[0]

		req := srch.GetViperParams()
		params, err := param.Parse(req.String())
		if err != nil {
			println(req.String())
			log.Fatal(err)
		}

		var docs []map[string]any

		var dataFile bool
		for _, r := range param.Routes {
			if viper.IsSet(r) {
				dataFile = true
			}
		}

		if dataFile {
			d := data.New(params.Route, params.Path)
			docs, err = d.Decode()
			if err != nil {
				println(req.String())
				log.Fatal(err)
			}
		} else {
			r := cmd.InOrStdin()
			err = data.DecodeNDJSON(r, &docs)
			if err != nil {
				println(req.String())
				log.Fatal(err)
			}
		}

		bi, err := srch.New(filepath.Join(param.Blv, path))
		if err != nil {
			log.Fatal(err)
		}

		for i, doc := range docs {
			id := cast.ToString(i)
			if di, ok := doc[params.UID]; ok {
				id = cast.ToString(di)
			}
			bi.Index(id, doc)
		}

	},
}

func init() {
	blvCmd.AddCommand(updateCmd)
}
