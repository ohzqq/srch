package srch

import (
	"encoding/json"
	"mime"
	"net/url"
	"path/filepath"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/sp"
	"github.com/ohzqq/srch/analyzer"
	"github.com/ohzqq/srch/doc"
	"github.com/spf13/cast"
)

type Idx struct {
	srch *hare.Table
	db   *hare.Database
	data *url.URL

	ID      int     `json:"_id"`
	Mapping Mapping `json:"mapping"`
	Name    string  `json:"name" qs:"name"`
	UID     string  `json:"uid,omitempty" mapstructure:"uid" qs:"uid"`

	// Index Settings
	SrchAttr  []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes,omitempty"`
	FacetAttr []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting,omitempty"`
	SortAttr  []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes,omitempty"`
}

func NewIdxCfg() *Idx {
	cfg := &Idx{
		SrchAttr: []string{"*"},
		Name:     "default",
	}
	return cfg.
		SetMapping(DefaultMapping())
}

func NewIdx() *Idx {
	return NewIdxCfg()
}

func (idx *Idx) ContentType() string {
	return mime.TypeByExtension(filepath.Ext(idx.data.Path))
}

func (idx *Idx) SetMapping(m Mapping) *Idx {
	idx.Mapping = m
	return idx
}

func (idx *Idx) SetDataURL(u *url.URL) *Idx {
	idx.data = u
	return idx
}

func (idx *Idx) idxTblName() string {
	return idx.Name + "Idx"
}

func (idx *Idx) dataTblName() string {
	return idx.Name + "Data"
}

func (idx *Idx) setSrchIdx(tbl *hare.Table) *Idx {
	idx.srch = tbl
	return idx
}

func (idx *Idx) setData(tbl *hare.Database) *Idx {
	idx.db = tbl
	return idx
}

func (idx *Idx) Decode(u url.Values) error {
	err := sp.Decode(u, idx)
	if err != nil {
		return err
	}
	idx.SrchAttr = parseSrchAttrs(idx.SrchAttr)
	if len(idx.FacetAttr) > 0 {
		idx.FacetAttr = ParseQueryStrings(idx.FacetAttr)
	}
	if len(idx.SortAttr) > 0 {
		idx.SortAttr = ParseQueryStrings(idx.SortAttr)
	}
	idx.SetMapping(idx.mapParams())
	return nil
}

func (idx *Idx) Cfg() *Idx {
	if idx.HasSrchAttr() || idx.HasFacetAttr() || idx.HasSortAttr() {
		return idx
	}
	return idx
}

func (idx *Idx) HasSrchAttr() bool {
	return len(idx.SrchAttr) > 0
}

func (idx *Idx) HasFacetAttr() bool {
	return len(idx.FacetAttr) > 0
}

func (idx *Idx) HasSortAttr() bool {
	return len(idx.SortAttr) > 0
}

func (idx *Idx) Encode() (url.Values, error) {
	return sp.Encode(idx)
}

func (idx *Idx) mapParams() Mapping {
	m := NewMapping()

	for _, attr := range idx.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range idx.FacetAttr {
		m.AddKeywords(attr)
	}

	for _, attr := range idx.SortAttr {
		m.AddKeywords(attr)
	}

	return m
}

func (idx *Idx) Insert(d []byte) error {
	doc := make(map[string]any)
	err := json.Unmarshal(d, &doc)
	if err != nil {
		return err
	}
	return idx.InsertDoc(doc)
}

func (idx *Idx) Find(ids ...int) ([]*doc.Doc, error) {
	var docs []*doc.Doc
	switch len(ids) {
	case 0:
		return docs, nil
	case 1:
		if ids[0] == -1 {
			return idx.FindAll()
		}
		fallthrough
	default:
		for _, id := range ids {
			doc := &doc.Doc{}
			err := idx.srch.Find(id, doc)
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
		return docs, nil
	}
}

func (idx *Idx) FindAll() ([]*doc.Doc, error) {
	ids, err := idx.srch.IDs()
	if err != nil {
		return nil, err
	}
	return idx.Find(ids...)
}

func (idx *Idx) InsertDoc(data map[string]any) error {
	doc := doc.New()
	for ana, attrs := range idx.Mapping {
		for field, val := range data {
			if ana == analyzer.Simple && slices.Equal(attrs, []string{"*"}) {
				doc.AddField(ana, field, val)
			}
			doc.AddField(ana, field, val)
		}
	}
	_, err := idx.srch.Insert(doc)
	if err != nil {
		return err
	}
	return nil
}

func (idx *Idx) Search(kw string) ([]int, error) {
	docs, err := idx.Find(-1)
	if err != nil {
		return nil, err
	}

	res := roaring.New()
	for ana, attrs := range idx.Mapping {
		for _, doc := range docs {
			for _, attr := range attrs {
				id := doc.Search(attr, ana, kw)
				if id != -1 {
					res.AddInt(id)
				}
			}
		}
	}
	ids := cast.ToIntSlice(res.ToArray())
	return ids, nil
}

func (idx *Idx) Changed(old *Idx) bool {
	if !slices.Equal(old.SrchAttr, idx.SrchAttr) {
		return true
	}
	if !slices.Equal(old.FacetAttr, idx.FacetAttr) {
		return true
	}
	if !slices.Equal(old.SortAttr, idx.SortAttr) {
		return true
	}
	if old.Name != idx.Name {
		return true
	}
	if old.UID != idx.UID {
		return true
	}
	return false
}

func NewCfgTbl(tbl string, m Mapping, id string) *Idx {
	return NewIdxCfg().
		SetMapping(m)
}

func DefaultIdxCfg() *Idx {
	return NewIdxCfg().
		SetMapping(DefaultMapping())
}

func (c *Idx) SetID(id int) {
	c.ID = id
}

func (c *Idx) GetID() int {
	return c.ID
}

func (c *Idx) AfterFind(db *hare.Database) error {
	//println("after find")
	return nil
}
