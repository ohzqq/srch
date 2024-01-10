package srch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type DataSrc func() []map[string]any

func SliceMapSrc(data []map[string]any) DataSrc {
	return func() []map[string]any {
		return data
	}
}

func StringSliceSrc(data []string) []map[string]any {
	d := make([]map[string]any, len(data))
	for i, item := range data {
		d[i] = map[string]any{"title": item}
	}
	return d
}

func FileSrc(file ...string) []map[string]any {
	data, err := NewDataFromFiles(file...)
	if err != nil {
		return []map[string]any{}
	}
	return data
}

func ReaderSrc(r io.Reader) []map[string]any {
	d, err := DecodeData(r)
	if err != nil {
		return []map[string]any{}
	}
	return d
}

// DecodeData decodes data from a io.Reader.
func DecodeData(r io.Reader) ([]map[string]any, error) {
	var data []map[string]any
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// NewDataFromFiles parses index data from files.
func NewDataFromFiles(d ...string) ([]map[string]any, error) {
	var data []map[string]any
	for _, datum := range d {
		p, err := dataFromFile(datum)
		if err != nil {
			return nil, err
		}
		data = append(data, p...)
	}
	return data, nil
}

func dataFromFile(d string) ([]map[string]any, error) {
	data, err := os.Open(d)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return DecodeData(data)
}

// NewDataFromString parses index data from a json formatted string.
func NewDataFromString(d string) ([]map[string]any, error) {
	buf := bytes.NewBufferString(d)
	return DecodeData(buf)
}

func parseData(d any) ([]map[string]any, error) {
	switch val := d.(type) {
	case []byte:
		return unmarshalData(val)
	case string:
		if exist(val) {
			return dataFromFile(val)
		} else {
			return unmarshalData([]byte(val))
		}
	case []map[string]any:
		return val, nil
	}
	return nil, errors.New("data couldn't be parsed")
}

func unmarshalData(d []byte) ([]map[string]any, error) {
	var data []map[string]any
	err := json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
