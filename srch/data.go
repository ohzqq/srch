package srch

import (
	"encoding/json"
	"io"

	"github.com/ohzqq/hare/dberr"
	"github.com/spf13/cast"
)

type FindItemFunc func(int) (map[string]any, error)

func (idx *Idx) findNdJSON(id int, r io.ReadCloser) (map[string]any, error) {
	dec := json.NewDecoder(r)
	i := 1
	for {
		item := make(map[string]any)
		if err := dec.Decode(&item); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		did := i
		if it, ok := item[idx.UID]; ok {
			did = cast.ToInt(it)
		}
		if did == id {
			return item, nil
		}
		i++
	}
	return nil, dberr.ErrNoRecord
}

func NdJSONFind(uid string, r io.ReadCloser) FindItemFunc {
	return func(id int) (map[string]any, error) {
		//r := bytes.NewReader(d)
		dec := json.NewDecoder(r)
		i := 1
		for {
			item := make(map[string]any)
			if err := dec.Decode(&item); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			did := i
			if it, ok := item[uid]; ok {
				did = cast.ToInt(it)
			}
			if did == id {
				return item, nil
			}
			i++
		}
		return nil, dberr.ErrNoRecord
	}
}
