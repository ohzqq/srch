package srch

import (
	"net/url"
)

const testValuesCfg = `and=tags:count:desc&field=title&or=authors:label:asc&or=narrators&or=series&data_file=testdata/data-dir/audiobooks.json&sort_by=title`

func queryParamsValues() url.Values {
	vals := make(url.Values)
	vals.Set("facetFilters", testComboFilter())
	return vals
}

func queryParamsString() string {
	return queryParamsValues().Encode()
}

func requestParams() string {
	p := queryParamsString()
	return p
}

func getNewQuery() url.Values {
	return ParseQuery(testValuesCfg, testQueryString, testSearchString)
}
