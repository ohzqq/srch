package cmd

import (
	"github.com/spf13/viper"
)

type flag int

//go:generate stringer -type flag -linecomment
const (
	B flag = iota // browse
	D             // dir
	F             // facet
	I             // index
	J             // json
	P             // params
	Q             // query
	R             // refine
	S             // search
	T             // text
	W             // workers
	U             // ui
)

func (f flag) Short() string {
	return string(f.String()[0])
}

func (f flag) Long() string {
	return f.String()
}

func init() {
	rootCmd.PersistentFlags().
		StringSliceP(
			I.Long(),
			I.Short(),
			[]string{},
			"list of data files to index",
		)
	rootCmd.PersistentFlags().
		StringP(
			D.Long(),
			D.Short(),
			"",
			"directory of data files",
		)
	rootCmd.PersistentFlags().
		StringP(
			J.Long(),
			J.Short(),
			"",
			"json formatted input",
		)

	//rootCmd.MarkFlagsOneRequired("index", "dir", "json")
	rootCmd.MarkFlagsMutuallyExclusive("index", "dir", "json")
	rootCmd.MarkPersistentFlagDirname("dir")
	rootCmd.MarkPersistentFlagFilename("index", ".json")

	rootCmd.PersistentFlags().
		Bool(
			U.Long(),
			false,
			"select results in a tui",
		)
	rootCmd.PersistentFlags().
		BoolP(
			B.Long(),
			B.Short(),
			false,
			"browse results in a tui",
		)

	rootCmd.PersistentFlags().
		StringSliceP(
			R.Long(),
			R.Short(),
			[]string{},
			"facet filters",
		)
	rootCmd.PersistentFlags().
		StringSliceP(
			S.Long(),
			S.Short(),
			[]string{},
			"text fields",
		)
	rootCmd.PersistentFlags().
		StringP(
			P.Long(),
			P.Short(),
			"",
			"encoded query/filter string (eg. color=red&color=pink&category=post",
		)
	rootCmd.PersistentFlags().
		StringP(
			Q.Long(),
			Q.Short(),
			"",
			"search index",
		)
	rootCmd.PersistentFlags().
		StringSliceP(
			"or",
			"o",
			[]string{},
			"disjunctive facets",
		)
	rootCmd.PersistentFlags().
		StringSliceP(
			"and",
			"a",
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
			W.Long(),
			W.Short(),
			1,
			"number of workers for computing facets",
		)
	viper.BindPFlag(
		W.Long(),
		rootCmd.Flags().Lookup("workers"),
	)

}
