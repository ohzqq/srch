package data

import (
	"strings"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Doc struct {
	SrchFields map[string]*bloom.BloomFilter `json:"searchableAttributes"`
	Facets     map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	id         int
}

func NewDoc(data map[string]any, params *param.Params) *Doc {
	doc := &Doc{
		SrchFields: make(map[string]*bloom.BloomFilter),
		Facets:     make(map[string]*bloom.BloomFilter),
		id:         cast.ToInt(data["id"]),
	}
	for _, attr := range params.SrchAttr {
		if f, ok := data[attr]; ok {
			str := cast.ToString(f)
			toks := facet.Tokenize(str)
			filter := bloom.NewWithEstimates(1000000, 0.01)
			for _, tok := range toks {
				filter.TestOrAddString(tok)
			}
			doc.SrchFields[attr] = filter
		}
	}

	for _, attr := range params.Facets {
		if f, ok := data[attr]; ok {
			str := cast.ToStringSlice(f)
			toks := facet.Keywords(str...)
			filter := bloom.NewWithEstimates(1000000, 0.01)
			for _, tok := range toks {
				filter.TestOrAddString(tok)
			}
			doc.Facets[attr] = filter
		}
	}
	return doc
}

func (d *Doc) SearchFields(kw string) []int {
	var ids []int
	for n, _ := range d.SrchFields {
		toks := facet.Tokenize(kw)
		var w []int
		for _, tok := range toks {
			if d.SearchField(n, tok) {
				w = append(w, d.id)
			}
		}
		if len(w) == len(toks) {
			ids = append(ids, d.id)
		}
	}
	return ids
}

func (d *Doc) SearchField(name string, kw string) bool {
	if f, ok := d.SrchFields[name]; ok {
		return f.TestString(kw)
	}
	return false
}

func (d *Doc) SearchFacets(kw string) bool {
	for _, f := range d.Facets {
		if ok := f.TestString(strings.ToLower(kw)); ok {
			return ok
		}
	}
	return false
}
