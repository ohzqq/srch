package srch

import (
	"encoding/json"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/hare"
	"github.com/ohzqq/sp"
	"github.com/spf13/cast"
)

type Idx struct {
	db      *hare.Database
	dataSrc *url.URL
	dataURL string
	idxURL  *url.URL

	ID         int     `json:"_id"`
	mapping    Mapping `json:"mapping"`
	Name       string  `json:"name" qs:"name"`
	PrimaryKey string  `json:"uid,omitempty" mapstructure:"uid" qs:"primaryKey"`

	// Index Settings
	SrchAttr  []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes,omitempty"`
	FacetAttr []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting,omitempty"`
	SortAttr  []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes,omitempty"`

	getData FindItemFunc
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

func (idx *Idx) Search(srch *Search) ([]map[string]any, error) {
	docs, err := idx.findDocs(nil)
	if err != nil {
		return nil, err
	}

	var ids []int
	if srch.Query == "" {
		ids = idx.getPKs(docs)
	} else {
		res := roaring.New()
		for ana, attrs := range idx.srchMap() {
			for _, doc := range docs {
				for _, attr := range attrs {
					id := doc.Search(attr, ana, srch.Query)
					if id != -1 {
						res.AddInt(id)
					}
				}
			}
		}
		ids = cast.ToIntSlice(res.ToArray())
	}

	data, err := idx.Find(ids...)
	if err != nil {
		return nil, err
	}
	data = srch.FilterRtrvAttr(data)
	return data, nil
}

func (idx *Idx) DataContentType() string {
	return mime.TypeByExtension(filepath.Ext(idx.dataSrc.Path))
}

func (idx *Idx) SetMapping(m Mapping) *Idx {
	idx.mapping = m
	return idx
}

func (idx *Idx) SetDataURL(u *url.URL) *Idx {
	idx.dataSrc = u
	return idx
}

func (idx *Idx) SetIdxURL(u *url.URL) *Idx {
	idx.idxURL = u
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
	return nil
}

func (idx *Idx) Cfg() *Idx {
	if idx.HasSrchAttr() || idx.HasFacetAttr() || idx.HasSortAttr() {
		return idx
	}
	return idx
}

func (idx *Idx) DocMapping() Mapping {
	if idx.mapping == nil ||
		idx.HasSrchAttr() ||
		idx.HasFacetAttr() ||
		idx.HasSortAttr() {
		return idx.mapParams()
	}
	return idx.mapping
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
	m := idx.srchMap()
	for _, attr := range idx.SortAttr {
		m.AddKeywords(attr)
	}
	return m
}

func (idx *Idx) srchMap() Mapping {
	m := NewMapping()

	for _, attr := range idx.SrchAttr {
		m.AddFulltext(attr)
	}

	for _, attr := range idx.FacetAttr {
		m.AddKeywords(attr)
	}

	return m
}

func (idx *Idx) AddDoc(d map[string]any) error {
	doc := New(d, idx.DocMapping(), idx.PrimaryKey)

	srch, err := idx.srch()
	if err != nil {
		return err
	}

	_, err = srch.Insert(doc)
	if err != nil {
		return err
	}

	return nil
}

func (idx *Idx) UpdateDoc(items ...map[string]any) error {
	pks := make([]int, len(items))
	for i, d := range items {
		pks[i] = getDocID(idx.PrimaryKey, d)
	}

	docs, err := idx.findDocByPK(pks...)
	if err != nil {
		return err
	}

	srch, err := idx.srch()
	if err != nil {
		return err
	}

	for _, doc := range docs {
		i := slices.Index(pks, doc.PrimaryKey)
		if i == -1 {
			continue
		}

		doc.Analyze(idx.DocMapping(), items[i])
		err = srch.Update(doc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (idx *Idx) Find(ids ...int) ([]map[string]any, error) {
	if idx.getData != nil {
		return idx.getData(ids...)
	}
	d, err := FindData(idx.dataSrc, ids)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (idx *Idx) findDocs(test func(*Doc) bool) ([]*Doc, error) {
	srch, err := idx.srch()
	if err != nil {
		return nil, err
	}

	ids, err := srch.IDs()
	if err != nil {
		return nil, err
	}

	var docs []*Doc
	for _, id := range ids {
		doc := DefaultDoc()
		err = srch.Find(id, doc)
		if err != nil {
			return nil, err
		}
		if test != nil {
			if test(doc) {
				docs = append(docs, doc)
			}
		} else {
			docs = append(docs, doc)
		}
	}
	return docs, nil
}

func (idx *Idx) getPKs(docs []*Doc) []int {
	pks := make([]int, len(docs))
	for i, doc := range docs {
		pks[i] = doc.PrimaryKey
	}
	return pks
}

func (idx *Idx) findDocByPK(pks ...int) ([]*Doc, error) {
	test := func(doc *Doc) bool {
		return slices.Contains(pks, doc.PrimaryKey)
	}
	return idx.findDocs(test)
}

func (idx *Idx) openData() (io.ReadCloser, error) {
	switch idx.dataSrc.Scheme {
	case "file":
		f, err := os.Open(idx.dataSrc.Path)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	return nil, nil
}

func (idx *Idx) Batch(r io.ReadCloser) error {
	dec := json.NewDecoder(r)
	for {
		d := make(map[string]any)
		if err := dec.Decode(&d); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err := idx.AddDoc(d)
		if err != nil {
			return err
		}
	}
	return nil
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
	if old.PrimaryKey != idx.PrimaryKey {
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

func (idx *Idx) data() (*hare.Table, error) {
	dataTbl := idx.Name + "Data"
	if !idx.db.TableExists(dataTbl) {
		err := idx.db.CreateTable(dataTbl)
		if err != nil {
			return nil, err
		}
	}
	tbl, err := idx.db.GetTable(dataTbl)
	if err != nil {
		return nil, err
	}
	return tbl, nil
}

func (idx *Idx) srch() (*hare.Table, error) {
	idxTbl := idx.Name + "Idx"
	if !idx.db.TableExists(idxTbl) {
		err := idx.db.CreateTable(idxTbl)
		if err != nil {
			return nil, err
		}
	}
	tbl, err := idx.db.GetTable(idxTbl)
	if err != nil {
		return nil, err
	}
	return tbl, nil
}

func (c *Idx) SetID(id int) {
	c.ID = id
}

func (c *Idx) GetID() int {
	return c.ID
}

func (idx *Idx) AfterFind(db *hare.Database) error {
	idx.db = db
	return nil
}
