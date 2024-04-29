package data

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyze"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Doc struct {
	Fields        map[string]*bloom.BloomFilter `json:"searchableAttributes"`
	Facets        map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	*param.Params `json:"-"`
	ID            int `json:"id"`
}

func NewDoc(data map[string]any, params *param.Params) *Doc {
	doc := &Doc{
		Fields: make(map[string]*bloom.BloomFilter),
		Facets: make(map[string]*bloom.BloomFilter),
		Params: params,
	}
	if id, ok := data[params.UID]; ok {
		doc.ID = cast.ToInt(id)
	}

	for _, attr := range params.SrchAttr {
		if f, ok := data[attr]; ok {
			str := cast.ToString(f)
			toks := analyze.Fulltext.Tokenize(str)
			filter := bloom.NewWithEstimates(uint(len(toks)*2), 0.01)
			for _, tok := range toks {
				filter.TestOrAddString(tok)
			}
			doc.Fields[attr] = filter
		}
	}

	for _, attr := range params.Facets {
		if f, ok := data[attr]; ok {
			str := cast.ToStringSlice(f)
			toks := analyze.Keywords.Tokenize(str...)
			filter := bloom.NewWithEstimates(uint(len(toks)*5), 0.01)
			for _, tok := range toks {
				filter.TestOrAddString(tok)
			}
			doc.Facets[attr] = filter
		}
	}
	return doc
}

func (d *Doc) SearchAllFields(kw string) bool {
	for n, _ := range d.Fields {
		return d.SearchField(n, kw)
	}
	return false
}

func (d *Doc) SearchField(name string, kw string) bool {
	if f, ok := d.Fields[name]; ok {
		toks := analyze.Fulltext.Tokenize(kw)
		for _, tok := range toks {
			if f.TestString(tok) {
				//fmt.Printf("id: %v, field: %v, token: %v\n", d.ID, name, tok)
				return true
			}
		}
	}
	return false
}

func (d *Doc) SearchFacets(kw string) []int {
	var ids []int
	for n, _ := range d.Fields {
		toks := analyze.Keywords.Tokenize(kw)
		var w []int
		for _, tok := range toks {
			if d.SearchFacet(n, tok) {
				w = append(w, d.ID)
			}
		}
		if len(w) == len(toks) {
			ids = append(ids, d.ID)
		}
	}
	return ids
}

func (d *Doc) SearchFacet(name string, kw string) bool {
	if f, ok := d.Facets[name]; ok {
		return f.TestString(kw)
	}
	return false
}

func (d *Doc) SetID(id int) {
	d.ID = id
}

func (d *Doc) GetID() int {
	return d.ID
}

func (d *Doc) AfterFind(_ *hare.Database) error {
	return nil
}
