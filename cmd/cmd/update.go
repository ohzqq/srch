package cmd

import (
	"log"

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
		var err error
		path := args[0]

		idx, docs, err := getIdxAndData(cmd, path)
		if err != nil {
			log.Fatal(err)
		}

		for i, doc := range docs {
			id := cast.ToString(i)
			if di, ok := doc[idx.Params.UID]; ok {
				id = cast.ToString(di)
			}
			err = idx.Index(id, doc)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	blvCmd.AddCommand(updateCmd)
}
