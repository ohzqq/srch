package cmd

import (
	"log"
	"path/filepath"

	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update INDEX",
	Short: "update documents in an index",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		req := srch.GetViperParams()

		idx, err := srch.New(req.String())
		if err != nil {
			println(req.String())
			log.Fatal(err)
		}

		bi, err := srch.New(filepath.Join(param.Blv, path))
		if err != nil {
			log.Fatal(err)
		}

		for i, doc := range idx.Docs {
			id := cast.ToString(i)
			if di, ok := doc[idx.Params.UID]; ok {
				id = cast.ToString(di)
			}
			bi.Index(id, doc)
		}

	},
}

func init() {
	blvCmd.AddCommand(updateCmd)
}
