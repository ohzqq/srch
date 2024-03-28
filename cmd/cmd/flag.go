package cmd

import (
	"github.com/ohzqq/srch/param"
	"github.com/spf13/viper"
)

func defineFlags() {
	for _, key := range param.SettingParams {
		short := string(key.String()[0])
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
		defineFlag(key.Slug(), short, usage)
	}

	for _, key := range param.SearchParams {
		short := string(key.String()[0])
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
		defineFlag(key.Slug(), short, usage)
	}

	for _, key := range param.Routes {
		short := string(key.String()[0])
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
		defineFlag(key.Slug(), short, usage)
	}

	defineFlag("and", "a", "conjunctive facets, format attribute:value")
	defineFlag("or", "o", "disjunctive facets, format attribute:value")
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
