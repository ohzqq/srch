package doc

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Doc struct {
	Standard map[string]*bloom.BloomFilter `json:"searchableAttributes"`
	Keyword  map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	Simple   map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	ID       int                           `json:"id"`
	Mapping  Mapping                       `json:"-"`
}

func New() *Doc {
	return &Doc{
		Standard: make(map[string]*bloom.BloomFilter),
		Keyword:  make(map[string]*bloom.BloomFilter),
		Simple:   make(map[string]*bloom.BloomFilter),
		Mapping:  NewMapping(),
	}
}

func (d *Doc) SetMapping(m Mapping) *Doc {
	d.Mapping = m
	return d
}

func (doc *Doc) AddField(ana analyzer.Analyzer, attr string, val any) {
	str := cast.ToString(val)
	toks := ana.Tokenize(str)
	filter := bloom.NewWithEstimates(uint(len(toks)*2), 0.01)
	for _, tok := range toks {
		filter.TestOrAddString(tok)
	}

	switch ana {
	case analyzer.Keyword:
		doc.Keyword[attr] = filter
	case analyzer.Standard:
		doc.Standard[attr] = filter
	case analyzer.Simple:
		fallthrough
	default:
		doc.Simple[attr] = filter
	}
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
					doc.Keyword[attr] = filter
				case analyzer.Standard:
					doc.Standard[attr] = filter
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
			doc.Standard[attr] = filter
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
			doc.Keyword[attr] = filter
		}
	}
	return doc
}

func (d *Doc) SearchAllFields(kw string) bool {
	for n, _ := range d.Standard {
		return d.SearchField(n, kw)
	}
	return false
}

func (d *Doc) Search(name string, ana analyzer.Analyzer, kw string) int {
	toks := ana.Tokenize(kw)
	var found []bool
	for _, tok := range toks {
		found = append(found, d.SearchField(name, tok))
	}
	if f := lo.Uniq(found); len(f) == 1 {
		if f[0] {
			return d.GetID()
		}
	}
	return -1
}

func (d *Doc) SearchField(name string, tok string) bool {
	if f, ok := d.Standard[name]; ok {
		if f.TestString(tok) {
			return true
		}
	}
	if f, ok := d.Keyword[name]; ok {
		if f.TestString(tok) {
			return true
		}
	}
	if f, ok := d.Simple[name]; ok {
		if f.TestString(tok) {
			return true
		}
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
