package cmd

import (
	"encoding/json"
	"log"
	"net/url"
	"path/filepath"

	"github.com/ohzqq/srch/param"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type flag int

type flags struct {
	set *pflag.FlagSet
}

//go:generate stringer -type flag -linecomment
const (
	A flag = iota // and
	B             // browse
	D             // dir
	F             // facet
	I             // index
	J             // json
	O             // or
	P             // params
	Q             // query
	R             // refine
	S             // search
	T             // fullText
	W             // workers
	U             // ui
)

func NewFlags(fl *pflag.FlagSet) *flags {
	return &flags{
		flags: fl,
	}
}

func GetViperRequest() *srch.Request {
	req := srch.NewRequest()
	for _, key := range param.SettingParams {
		switch key {
		case param.SrchAttr:
			val := viper.GetStringSlice(key)
			req.SrchAttr(val...)
		case param.FacetAttr:
			val := viper.GetStringSlice(key)
			req.SrchAttr(val...)
		case param.SortAttr:
			val := viper.GetStringSlice(key)
			req.SrchAttr(val...)
		case param.UID:
			val := viper.GetString(key)
			req.UID(val)
		case param.Format:
			val := viper.GetString(key)
			req.Format(val)
		}
	}

	for _, key := range param.SearchParams {
		switch key {
		case param.Route:
			val := viper.GetString(key)
			req.SetRoute(val)
		case param.SortFacetsBy:
			val := viper.GetString(key)
			req.SortFacetsBy(val)
		case param.Facets:
			val := viper.GetStringSlice(key)
			req.Facets(val...)
		case param.Filters:
			val := viper.GetString(key)
			req.Filters(val)
		case "or":
			val := viper.GetStringSlice("or")
			req.OrFilter(val...)
		case "and":
			val := viper.GetStringSlice("and")
			req.AndFilter(val...)
		case param.RtrvAttr:
			val := viper.GetStringSlice(key)
			req.RtrvAttr(val...)
		case param.Page:
			val := viper.GetInt(key)
			req.Page(val)
		case param.HitsPerPage:
			val := viper.GetInt(key)
			req.HitsPerPage(val)
		case param.Query:
			val := viper.GetString(key)
			req.Query(val)
		case param.SortBy:
			val := viper.GetString(key)
			req.SortBy(val)
		case param.Order:
			val := viper.GetString(key)
			req.Order(val)
		}
	}

	for _, key := range param.Routes {
		switch key {
		case param.Blv:
			val := viper.GetString(key)
			req.SetRoute(filepath.Join(key, val))
		case param.Dir:
			val := viper.GetString(key)
			req.SetRoute(filepath.Join(key, val))
		case param.File:
			val := viper.GetStringSlice(key)
			req.SetRoute(filepath.Join(key, val))
		}
	}
	return req
}

func FlagsToRequest(flags *pflag.FlagSet) string {
}

func (f flag) Short() string {
	return string(f.String()[0])
}

func (f flag) Long() string {
	return f.String()
}

func (f flag) Param() string {
	switch f {
	case A:
	case O:
	case B:
	case D:
		return param.DataDir
	case F:
		return param.FacetAttr
	case I:
		return param.DataFile
	case J:
	case P:
	case Q:
		return param.Query
	case R:
	case S:
		return param.SrchAttr
	case T:
		return param.FullText
	case W:
	case U:
	}
	return ""
}

var allFlags = []flag{
	A,
	B,
	D,
	F,
	I,
	J,
	O,
	P,
	Q,
	R,
	S,
	T,
	W,
	U,
}

func (f flag) GetSlice(flags *pflag.FlagSet) []string {
	and, err := flags.GetStringSlice(f.Long())
	if err != nil {
		return []string{}
	}
	return and
}

func (f flag) GetString(flags *pflag.FlagSet) string {
	and, err := flags.GetString(f.Long())
	if err != nil {
		return ""
	}
	return and
}

func getSlice(flags *pflag.FlagSet, long string) []string {
	and, err := flags.GetStringSlice(long)
	if err != nil {
		log.Fatal(err)
	}
	return and
}

func getString(flags *pflag.FlagSet, long string) string {
	and, err := flags.GetString(long)
	if err != nil {
		log.Fatal(err)
	}
	return and
}

func FlagsToParams(flags *pflag.FlagSet) url.Values {
	params := make(url.Values)
	var filters []any
	for _, flag := range allFlags {
		long := flag.Long()
		param := flag.Param()

		if !flags.Changed(long) {
			continue
		}

		switch flag {
		case A:
			for _, a := range flag.GetSlice(flags) {
				filters = append(filters, a)
			}
		case O:
			filters = append(filters, flag.GetSlice(flags))
		case F, I, S:
			params[param] = flag.GetSlice(flags)
		case P:
			params = srch.ParseQuery(flag.GetString(flags))
		case Q, D, T:
			params.Set(param, flag.GetString(flags))
		}
	}

	if len(filters) > 0 {
		d, err := json.Marshal(filters)
		if err != nil {
			log.Fatal(err)
		}
		params.Set(srch.FacetFilters, string(d))
	}

	return params
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
			F.Long(),
			F.Short(),
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
			O.Long(),
			O.Short(),
			[]string{},
			"disjunctive facets",
		)
	rootCmd.PersistentFlags().
		StringSliceP(
			A.Long(),
			A.Short(),
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
