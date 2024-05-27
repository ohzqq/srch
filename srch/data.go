package srch

import (
	"encoding/json"
	"io"
	"slices"

	"github.com/ohzqq/hare/dberr"
	"github.com/spf13/cast"
)

type FindItemFunc func(...int) ([]map[string]any, error)

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
		if it, ok := item[idx.PrimaryKey]; ok {
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
	return func(ids ...int) ([]map[string]any, error) {
		//r := bytes.NewReader(d)
		dec := json.NewDecoder(r)
		i := 1
		var items []map[string]any
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
			if slices.Contains(ids, did) {
				items = append(items, item)
			}
			i++
		}
		return items, nil
	}
}

func findNDJSON(r io.Reader, uid string, ids ...int) ([]map[string]any, error) {
	dec := json.NewDecoder(r)
	i := 1
	var items []map[string]any
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
		if len(ids) > 1 {
			if slices.Contains(ids, did) {
				items = append(items, item)
			}
		} else {
			items = append(items, item)
		}
		i++
	}
	return items, nil
}
