package srch

import (
	"encoding/json"
	"io"
	"slices"

	"github.com/spf13/cast"
)

type FindItemFunc func(...int) ([]map[string]any, error)

func SrcNDJSON(r io.ReadCloser, pk string) FindItemFunc {
	return func(ids ...int) ([]map[string]any, error) {
		items, err := findNDJSON(r, pk, ids...)
		if err != nil {
			return nil, err
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
		if len(ids) > 0 {
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
