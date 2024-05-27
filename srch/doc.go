package srch

import (
	"slices"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/srch/analyzer"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Doc struct {
	Standard map[string]*bloom.BloomFilter `json:"searchableAttributes"`
	Keyword  map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	Simple   map[string]*bloom.BloomFilter `json:"attributesForFaceting"`
	ID       int                           `json:"_id"`
	UID      int                           `json:"id,omitempty"`
	CustomID string                        `json:"-"`
}

func New(key ...string) *Doc {
	var pk string
	if len(key) > 0 {
		pk = key[0]
	}
	return DefaultDoc().
		WithCustomID(pk)
}

func DefaultDoc() *Doc {
	return &Doc{
		Standard: make(map[string]*bloom.BloomFilter),
		Keyword:  make(map[string]*bloom.BloomFilter),
		Simple:   make(map[string]*bloom.BloomFilter),
	}
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

func (doc *Doc) WithCustomID(f string) *Doc {
	doc.CustomID = f
	return doc
}

func (doc *Doc) SetData(m Mapping, data map[string]any) *Doc {
	for ana, attrs := range m {
		for field, val := range data {
			for _, attr := range attrs {
				if field == attr {
					if ana == analyzer.Simple && slices.Equal(attrs, []string{"*"}) {
						doc.AddField(ana, field, val)
					}
					doc.AddField(ana, field, val)
				}
			}
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
		f := d.SearchField(name, tok)
		found = append(found, f)
	}

	if f := lo.Uniq(found); len(f) == 1 {
		if f[0] {
			//fmt.Printf("field %s: found %v\n", name, found)
			return d.UID
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
	if d.UID < 1 {
		d.UID = id
	}
}

func (d *Doc) GetID() int {
	return d.ID
}

func (d *Doc) AfterFind(db *hare.Database) error {
	return nil
}

func getDocID(uid any, doc map[string]any) int {
	if u, ok := doc[cast.ToString(uid)]; ok {
		return cast.ToInt(u)
	}
	return cast.ToInt(uid)
}
