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
		long := flect.Dasherize(key)
		short := string(key[0])
		var usage string
		switch key {
		case param.UID:
			usage = `uid for documents in collection, default: id`
		case param.SrchAttr:
			short = "t"
			usage = `list of fields to search, use the special label '*' to search all fields, default: all`
		default:
			continue
		}
		defineFlag(long, short, usage)
	}

	for _, key := range param.SearchParams {
		long := flect.Dasherize(key)
		short := string(key[0])
		var usage string
		switch key {
		case param.Facets:
			usage = `list of facets, format is attribute[:sort][:order], eg 'tags:count:desc'`
		case param.RtrvAttr:
			short = ""
			usage = `list of fields to retrieve for display, default: all`
		case param.Page:
			usage = `page number for paginated results, default: 0`
		case param.Query:
			usage = `query for search`
		case param.HitsPerPage:
			short = "l"
			usage = `number of hits to return, default: all`
		case param.SortBy:
			usage = `field to sort results by, default: uid:desc`
		default:
			continue
		}
		defineFlag(long, short, usage)
	}

	for _, key := range param.Routes {
		long := flect.Dasherize(key)
		short := string(key[0])
		var usage string
		switch key {
		case param.File:
			short = "j"
			usage = "json or ndjson file to index"
		case param.Dir:
			usage = "directory of data files"
		case param.Blv:
			usage = "path to bleve index"
		}
		defineFlag(long, short, usage)
	}

	defineFlag("and", "a", "conjuctive facets, format attribute:value")
	defineFlag("or", "o", "disconjuctive facets, format attribute:value")
}

func defineFlag(long, short, usage string) {
	if short != "" {
		rootCmd.PersistentFlags().
			StringSliceP(
				long,
				short,
				[]string{},
				usage,
			)
	} else {
		rootCmd.PersistentFlags().
			StringSlice(
				long,
				[]string{},
				usage,
			)
	}
	viper.BindPFlag(
		long,
		rootCmd.PersistentFlags().Lookup(long),
	)
}
