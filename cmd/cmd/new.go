package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/blv"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new DIR",
	Short: `Create a new bleve index`,
	Long:  `Create a new bleve index`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		route := filepath.Join(param.Blv.String(), args[0])
		req := srch.NewRequest().SetRoute(route)

		params, err := param.Parse(req.String())
		if err != nil {
			log.Fatal(err)
		}

		idx, err := blv.New(params)
		if err != nil {
			log.Fatalf("%v at %s\n", err, idx.Path)
		}

		fmt.Printf("new index created at %s\n", idx.Path)
	},
}

func init() {
	blvCmd.AddCommand(newCmd)
}
