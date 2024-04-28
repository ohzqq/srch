package data

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/srch/facet"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Doc struct {
	SrchFields map[string]*bloom.BloomFilter `json:"searchableAttributes"`
	Facets     map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
}

func NewDoc(data map[string]any, params *param.Params) *Doc {
	doc := &Doc{
		SrchFields: make(map[string]*bloom.BloomFilter),
		Facets:     make(map[string]*bloom.BloomFilter),
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
	return doc
}

func (d *Doc) Test(kw string) bool {
	for _, f := range d.SrchFields {
		if ok := f.TestString(kw); ok {
			return ok
		}
	}
	return false
}
