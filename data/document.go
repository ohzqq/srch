package data

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyze"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Doc struct {
	Fulltext      map[string]*bloom.BloomFilter `json:"searchableAttributes"`
	Keyword       map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	*param.Params `json:"-"`
	ID            int `json:"id"`
}

func newDoc() *Doc {
	return &Doc{
		Fulltext: make(map[string]*bloom.BloomFilter),
		Keyword:  make(map[string]*bloom.BloomFilter),
	}
}

func NewDoc(data map[string]any, params *param.Params) *Doc {
	doc := newDoc()
	for _, attr := range params.SrchAttr {
		if f, ok := data[attr]; ok {
			str := cast.ToString(f)
			toks := analyze.Fulltext.Tokenize(str)
			filter := bloom.NewWithEstimates(uint(len(toks)*2), 0.01)
			for _, tok := range toks {
				filter.TestOrAddString(tok)
			}
			doc.Fulltext[attr] = filter
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
			doc.Keyword[attr] = filter
		}
	}
	return doc
}

func (d *Doc) SearchAllFields(kw string) bool {
	for n, _ := range d.Fulltext {
		return d.SearchField(n, kw)
	}
	return false
}

func (d *Doc) SearchField(name string, kw string) bool {
	if f, ok := d.Fulltext[name]; ok {
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
	for n, _ := range d.Fulltext {
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
	if f, ok := d.Keyword[name]; ok {
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

func getDocID(uid any, doc map[string]any) int {
	if u, ok := doc[cast.ToString(uid)]; ok {
		return cast.ToInt(u)
	}
	return cast.ToInt(uid)
}
