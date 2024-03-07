package fuzzy

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/ohzqq/srch"
	"github.com/spf13/cast"
)

type Indexer struct {
	fields  []string
	data    []map[string]any
	bits    *roaring.Bitmap
	*Params `json:"params"`
}

func New(settings string) (srch.Searcher, error) {
	idx := &Indexer{
		bits: roaring.New(),
	}

	return idx.Open(settings)
}

func (idx *Indexer) Open(settings string) (srch.Searcher, error) {
	idx.Params = srch.ParseParams(settings)

	if !idx.Params.HasData() {
		return nil, NoDataErr
	}

	var err error
	var data []map[string]any
	switch {
	case idx.Params.Has(DataFile):
		data, err = srch.FileSrc(idx.GetSlice(DataFile)...)
		idx.Settings.Del(DataFile)
	case idx.Has(DataDir):
		data, err = srch.DirSrc(idx.Get(DataDir))
		idx.Settings.Del(DataDir)
	}
	if err != nil {
		return idx, err
	}

	err = idx.Index("", data...)

	return idx, nil
}

func (idx *Indexer) Index(uid string, data ...[]map[string]any) error {
	idx.data = data

	for id, d := range idx.data {
		idx.bits.AddInt(parseID(id, idx.Params.UID(), d))
	}

	return nil
}

// String satisfies the fuzzy.Source interface.
func (idx *Indexer) String(i int) string {
	for _, d := range idx.data {
		id := parseID(i, idx.Params.UID(), d)
	}
	return idx.data[i]
}

// Len satisfies the fuzzy.Source interface.
func (idx *Indexer) Len() int {
	return len(idx.data)
}

func parseID(id int, uid string, d map[string]any) int {
	if uid == "" {
		return id
	}
	if i, ok := d[uid]; ok {
		return cast.ToInt(i)
	}
	return id
}

func parseSearchableFields(attr []string, d map[string]any) string {

	var str string
	for _, a := range attr {
		if v, ok := d[a]; ok {
			str += cast.ToString(v)
			str += " "
		}
	}
}
