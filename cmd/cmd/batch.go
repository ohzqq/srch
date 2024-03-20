package cmd

import (
	"log"

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

		idx, docs, err := getIdxAndData(cmd, path)
		if err != nil {
			log.Fatal(err)
		}

		err = idx.Batch(docs)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	blvCmd.AddCommand(batchCmd)
}
