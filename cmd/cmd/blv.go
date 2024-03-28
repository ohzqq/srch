package cmd

import (
	"fmt"
	"net/url"

	"github.com/ohzqq/srch"
	"github.com/ohzqq/srch/data"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// blvCmd represents the blv command
var blvCmd = &cobra.Command{
	Use:   "blv",
	Short: "work with blv indexes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("blv called")
	},
}

func getIdxAndData(cmd *cobra.Command, path string) (*srch.Index, []map[string]any, error) {
	req := srch.GetViperParams()
	params, err := param.Parse(req.String())
	if err != nil {
		return nil, nil, fmt.Errorf("%s failed to parse err: %w\n", err)
	}

	docs, err := getData(cmd, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get data err: %w\n", err)
	}

	u := &url.URL{
		Path:     param.Blv.String(),
		RawQuery: "path=" + path,
	}

	idx, err := srch.New(u.String())
	return idx, docs, err
}

func getData(cmd *cobra.Command, params *param.Params) ([]map[string]any, error) {
	var docs []map[string]any
	var err error

	var dataFile bool
	for _, r := range param.Routes {
		if viper.IsSet(r.Snake()) {
			dataFile = true
		}
	}

	if dataFile {
		d := data.New(params.Route, params.Path)
		docs, err = d.Decode()
		if err != nil {
			return nil, fmt.Errorf("failed to decode at %s data err: %w\n", params.Path, err)
		}
	} else {
		r := cmd.InOrStdin()
		err = data.DecodeNDJSON(r, &docs)
		if err != nil {
			return nil, fmt.Errorf("failed to decode data from stdin err: %w\n", err)
		}
	}

	return docs, nil
}

func init() {
	rootCmd.AddCommand(blvCmd)
}
