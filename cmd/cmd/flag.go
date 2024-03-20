package cmd

import (
	"github.com/gobuffalo/flect"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/viper"
)

type flag int

//go:generate stringer -type flag -linecomment
const (
	And     flag = iota // and
	Blv                 // browse
	Dir                 // dir
	Facets              // facet
	Index               // index
	JSON                // json
	Or                  // or
	Params              // params
	Query               // query
	Refine              // refine
	Search              // search
	Workers             // workers
	UI                  // ui
)

func (f flag) Short() string {
	return string(f.String()[0])
}

func (f flag) Long() string {
	return f.String()
}

func (f flag) Param() string {
	switch f {
	case And:
	case Or:
	case Blv:
	case Dir:
		return param.DataDir
	case Facets:
		return param.FacetAttr
	case Index:
		return param.DataFile
	case JSON:
	case Params:
	case Query:
		return param.Query
	case Refine:
	case Search:
		return param.SrchAttr
	case Workers:
	case UI:
	}
	return ""
}

var allFlags = []flag{
	And,
	Blv,
	Dir,
	Facets,
	Index,
	JSON,
	Or,
	Params,
	Query,
	Refine,
	Search,
	Workers,
	UI,
}

func defineFlags() {
	for _, key := range param.SettingParams {
		long := flect.New(key).Dasherize()
		short := key[0]
		switch key {
		case param.UID:
		case param.Format:
		default:
			rootCmd.PersistentFlags().
				StringSliceP(
					long,
					short,
					[]string{},
					"list of data files to index",
				)
		}
	}
	for _, key := range param.SearchParams {
		switch key {
		}
	}
	for _, key := range param.Routes {
		switch key {
		}
	}
}

func init() {
	rootCmd.PersistentFlags().
		StringSliceP(
			Index.Long(),
			Index.Short(),
			[]string{},
			"list of data files to index",
		)
	rootCmd.PersistentFlags().
		StringP(
			Dir.Long(),
			Dir.Short(),
			"",
			"directory of data files",
		)
	rootCmd.PersistentFlags().
		StringP(
			JSON.Long(),
			JSON.Short(),
			"",
			"json formatted input",
		)

	//rootCmd.MarkFlagsOneRequired("index", "dir", "json")
	rootCmd.MarkFlagsMutuallyExclusive("index", "dir", "json")
	rootCmd.MarkPersistentFlagDirname("dir")
	rootCmd.MarkPersistentFlagFilename("index", ".json")

	rootCmd.PersistentFlags().
		Bool(
			UI.Long(),
			false,
			"select results in a tui",
		)
	rootCmd.PersistentFlags().
		BoolP(
			Blv.Long(),
			Blv.Short(),
			false,
			"browse results in a tui",
		)

	rootCmd.PersistentFlags().
		StringSliceP(
			Facets.Long(),
			Facets.Short(),
			[]string{},
			"facet filters",
		)
	rootCmd.PersistentFlags().
		StringSliceP(
			Search.Long(),
			Search.Short(),
			[]string{},
			"text fields",
		)
	rootCmd.PersistentFlags().
		StringP(
			Params.Long(),
			Params.Short(),
			"",
			"encoded query/filter string (eg. color=red&color=pink&category=post",
		)
	rootCmd.PersistentFlags().
		StringP(
			Query.Long(),
			Query.Short(),
			"",
			"search index",
		)
	rootCmd.PersistentFlags().
		StringSliceP(
			Or.Long(),
			Or.Short(),
			[]string{},
			"disjunctive facets",
		)
	rootCmd.PersistentFlags().
		StringSliceP(
			And.Long(),
			And.Short(),
			[]string{},
			"conjunctive facets",
		)

	rootCmd.PersistentFlags().
		Bool(
			"pretty",
			false,
			"pretty print json output",
		)

	rootCmd.PersistentFlags().
		IntP(
			Workers.Long(),
			Workers.Short(),
			1,
			"number of workers for computing facets",
		)
	viper.BindPFlag(
		Workers.Long(),
		rootCmd.Flags().Lookup("workers"),
	)

}
