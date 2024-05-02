package doc

import (
	"fmt"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

type Doc struct {
	Fulltext map[string]*bloom.BloomFilter `json:"searchableAttributes"`
	Keywords map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	Simple   map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	ID       int                           `json:"id"`
	Mapping  Mapping                       `json:"-"`
}

func New() *Doc {
	return &Doc{
		Fulltext: make(map[string]*bloom.BloomFilter),
		Keywords: make(map[string]*bloom.BloomFilter),
		Simple:   make(map[string]*bloom.BloomFilter),
		Mapping:  NewMapping(),
	}
}

func (d *Doc) SetMapping(m Mapping) *Doc {
	d.Mapping = m
	return d
}

func (doc *Doc) SetData(data map[string]any) *Doc {
	for ana, attrs := range doc.Mapping {
		for _, attr := range attrs {
			if val, ok := data[attr]; ok {
				str := cast.ToString(val)
				toks := ana.Tokenize(str)
				filter := bloom.NewWithEstimates(uint(len(toks)*2), 0.01)
				for _, tok := range toks {
					filter.TestOrAddString(tok)
				}

				switch ana {
				case analyzer.Keyword:
					doc.Keywords[attr] = filter
				case analyzer.Standard:
					doc.Fulltext[attr] = filter
				case analyzer.Simple:
					fallthrough
				default:
					doc.Simple[attr] = filter
				}
			}
		}
	}
	return doc
}

func NewDoc(data map[string]any, params *param.Params) *Doc {
	doc := New()
	for _, attr := range params.SrchAttr {
		if f, ok := data[attr]; ok {
			str := cast.ToString(f)
			toks := analyzer.Standard.Tokenize(str)
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
			toks := analyzer.Keyword.Tokenize(str...)
			filter := bloom.NewWithEstimates(uint(len(toks)*5), 0.01)
			for _, tok := range toks {
				filter.TestOrAddString(tok)
			}
			doc.Keywords[attr] = filter
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

func (d *Doc) Search(name string, ana analyzer.Analyzer, kw string) int {
	toks := ana.Tokenize(kw)
	switch ana {
	case analyzer.Standard:
		if f, ok := d.Fulltext[name]; ok {
			for _, tok := range toks {
				if f.TestString(tok) {
					fmt.Printf("id: %v, field: %v, token: %v\n", d.ID, name, tok)
					return d.GetID()
				}
			}
		}
	case analyzer.Keyword:
		if f, ok := d.Keywords[name]; ok {
			for _, tok := range toks {
				if f.TestString(tok) {
					//fmt.Printf("id: %v, field: %v, token: %v\n", d.ID, name, tok)
					return d.GetID()
				}
			}
		}
	case analyzer.Simple:
	}
	return -1
}

func (d *Doc) SearchField(name string, kw string) bool {
	if f, ok := d.Fulltext[name]; ok {
		toks := analyzer.Standard.Tokenize(kw)
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
	for n, _ := range d.Keywords {
		toks := analyzer.Keyword.Tokenize(kw)
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
	if f, ok := d.Keywords[name]; ok {
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
