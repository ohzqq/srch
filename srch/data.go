package srch

import (
	"encoding/json"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cast"
)

const (
	NdJSON = `application/x-ndjson`
	JSON   = `application/json`
	Hare   = `application/hare`
)

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
	mime.AddExtensionType(".hare", "application/hare")
}

type FindItemFunc func([]int) []map[string]any
type FindItems[T any] func(...T) ([]map[string]any, error)

func FindData[T any](u *url.URL, ids []T) []map[string]any {
	ct := mime.TypeByExtension(filepath.Ext(u.Path))
	switch ct {
	case NdJSON:
		return SrcNDJSON(u, ids)
	}

	return []map[string]any{}
}

func SrcNDJSON[T any](u *url.URL, ids []T) []map[string]any {
	var err error
	var r io.ReadCloser
	switch u.Scheme {
	case "file":
		r, err = os.Open(u.Path)
		if err != nil {
			return []map[string]any{}
		}
	case "http", "https":
		res, err := client.Get(u.String())
		if err != nil {
			return []map[string]any{}
		}
		r = res.Body
	default:
		return []map[string]any{}
	}
	defer r.Close()
	return findNDJSON(r, u.Query().Get("primaryKey"), ids)
}

func findNDJSON[T any](r io.Reader, uid string, ids []T) []map[string]any {
	dec := json.NewDecoder(r)
	i := 1
	guids := cast.ToIntSlice(ids)
	var items []map[string]any
	for {
		item := make(map[string]any)
		if err := dec.Decode(&item); err == io.EOF {
			break
		} else if err != nil {
			return items
		}
		did := i
		if it, ok := item[uid]; ok {
			did = cast.ToInt(it)
		}
		if len(ids) > 0 {
			if slices.Contains(guids, did) {
				items = append(items, item)
			}
		} else {
			items = append(items, item)
		}
		i++
	}
	return items
}
