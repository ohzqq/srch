package srch

import (
	"encoding/json"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"github.com/ohzqq/hare"
	"github.com/ohzqq/sp"
	"github.com/spf13/cast"
)

type Idx struct {
	db      *hare.Database
	dataURL *url.URL
	idxURL  *url.URL

	ID      int     `json:"_id"`
	Mapping Mapping `json:"mapping"`
	Name    string  `json:"name" qs:"name"`
	UID     string  `json:"uid,omitempty" mapstructure:"uid" qs:"uid"`

	// Index Settings
	SrchAttr  []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes,omitempty"`
	FacetAttr []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting,omitempty"`
	SortAttr  []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes,omitempty"`

	getData FindItemFunc
}

const (
	NdJSON = `application/x-ndjson`
	JSON   = `application/json`
	Hare   = `application/hare`
)

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

func (idx *Idx) DataContentType() string {
	return mime.TypeByExtension(filepath.Ext(idx.dataURL.Path))
}

func (idx *Idx) SetMapping(m Mapping) *Idx {
	idx.Mapping = m
	return idx
}

func (idx *Idx) SetDataURL(u *url.URL) *Idx {
	idx.dataURL = u
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

func (idx *Idx) Find(id int) (map[string]any, error) {
	doc := New()
	srch, err := idx.srch()
	if err != nil {
		return nil, err
	}
	err = srch.Find(id, doc)
	if err != nil {
		return nil, err
	}

	d, err := idx.getData(doc.UID)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (idx *Idx) Insert(item *Item) error {
	doc := item.Idx(idx.Mapping).
		WithCustomID(idx.UID)

	if id, ok := item.Data[idx.UID]; ok {
		i := cast.ToInt(id)
		doc.UID = i
	}

	srch, err := idx.srch()
	if err != nil {
		return err
	}

	_, err = srch.Insert(doc)
	if err != nil {
		return err
	}

	//data, err := idx.data()
	//if err != nil {
	//return err
	//}

	//_, err = data.Insert(item)
	//if err != nil {
	//return err
	//}

	return nil
}

func (idx *Idx) openData() (io.ReadCloser, error) {
	switch idx.dataURL.Scheme {
	case "file":
		f, err := os.Open(idx.dataURL.Path)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	return nil, nil
}

func (idx *Idx) Batch(r io.ReadCloser) error {
	//r := bytes.NewReader(d)
	dec := json.NewDecoder(r)
	for {
		item := NewItem()
		if err := dec.Decode(&item.Data); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err := idx.Insert(item)
		if err != nil {
			return err
		}
	}
	return nil
}

//func (idx *Idx) Insert(d []byte) error {
//  doc := make(map[string]any)
//  err := json.Unmarshal(d, &doc)
//  if err != nil {
//    return err
//  }
//  return idx.InsertDoc(doc)
//}

//func (idx *Idx) Find(ids ...int) ([]*doc.Doc, error) {
//  var docs []*doc.Doc
//  switch len(ids) {
//  case 0:
//    return docs, nil
//  case 1:
//    if ids[0] == -1 {
//      return idx.FindAll()
//    }
//    fallthrough
//  default:
//    for _, id := range ids {
//      doc := &doc.Doc{}
//      err := idx.srch.Find(id, doc)
//      if err != nil {
//        return nil, err
//      }
//      docs = append(docs, doc)
//    }
//    return docs, nil
//  }
//}

//func (idx *Idx) FindAll() ([]*doc.Doc, error) {
//  ids, err := idx.srch.IDs()
//  if err != nil {
//    return nil, err
//  }
//  return idx.Find(ids...)
//}

//func (idx *Idx) InsertDoc(data map[string]any) error {
//  doc := doc.New()
//  for ana, attrs := range idx.Mapping {
//    for field, val := range data {
//      if ana == analyzer.Simple && slices.Equal(attrs, []string{"*"}) {
//        doc.AddField(ana, field, val)
//      }
//      doc.AddField(ana, field, val)
//    }
//  }
//  _, err := idx.srch.Insert(doc)
//  if err != nil {
//    return err
//  }
//  return nil
//}

//func (idx *Idx) Search(kw string) ([]int, error) {
//  docs, err := idx.Find(-1)
//  if err != nil {
//    return nil, err
//  }

//  res := roaring.New()
//  for ana, attrs := range idx.Mapping {
//    for _, doc := range docs {
//      for _, attr := range attrs {
//        id := doc.Search(attr, ana, kw)
//        if id != -1 {
//          res.AddInt(id)
//        }
//      }
//    }
//  }
//  ids := cast.ToIntSlice(res.ToArray())
//  return ids, nil
//}

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

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
	mime.AddExtensionType(".hare", "application/hare")
}
