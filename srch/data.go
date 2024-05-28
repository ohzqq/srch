package srch

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"slices"

	"github.com/spf13/cast"
)

type DataSrc struct {
	*url.URL
}

func NewDataSrc(uri string) *DataSrc {
	var err error
	d := &DataSrc{}
	d.URL, err = parseURL(uri)
	if err != nil {
		d.URL = &url.URL{Scheme: "mem"}
	}
	return d
}

func (d *DataSrc) Find(ids []any) ([]map[string]any, error) {
	return nil, nil
}

func FindData[T any](q string, col []T, fn func(uri *url.URL, ids ...T) ([]map[string]any, error)) ([]map[string]any, error) {
	u, err := parseURL(q)
	if err != nil {
		u = &url.URL{Scheme: "mem"}
	}
	return fn(u, col...)
}

func SrcNDJSON[T any](u *url.URL, ids ...T) ([]map[string]any, error) {
	var err error
	var r io.ReadCloser
	switch u.Scheme {
	case "file":
		r, err = os.Open(u.Path)
		if err != nil {
			return nil, err
		}
	case "http", "https":
		res, err := client.Get(u.String())
		if err != nil {
			return nil, err
		}
		r = res.Body
	}
	defer r.Close()
	return findNDJSONz(r, u.Query().Get("primaryKey"), ids...)
}

type FindItemFunc func(...int) ([]map[string]any, error)
type FindItems[T any] func(...T) ([]map[string]any, error)

func NDJSONsrc(r io.ReadCloser, pk string) FindItemFunc {
	return func(ids ...int) ([]map[string]any, error) {
		items, err := findNDJSON(r, pk, ids...)
		if err != nil {
			return nil, err
		}

		return items, nil
	}
}

func findNDJSONz[T any](r io.Reader, uid string, ids ...T) ([]map[string]any, error) {
	dec := json.NewDecoder(r)
	i := 1
	guids := cast.ToIntSlice(ids)
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
			if slices.Contains(guids, did) {
				items = append(items, item)
			}
		} else {
			items = append(items, item)
		}
		i++
	}
	return items, nil
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
