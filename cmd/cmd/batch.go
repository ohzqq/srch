package cmd

import (
	"log"
	"path/filepath"

	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cobra"
)

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch INDEX",
	Short: "batch add documents to index",
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

		err = bi.Batch(idx.Docs)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	blvCmd.AddCommand(batchCmd)
}
