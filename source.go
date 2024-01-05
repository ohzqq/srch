package srch

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Src func(args ...any) []any

func SliceSrc(data ...any) Src {
	return func(...any) []any {
		return data
	}
}

func FileSrc(file ...string) Src {
	return func(...any) []any {
		data, err := NewDataFromFiles(file...)
		if err != nil {
			return []any{}
		}
		return data
	}
}

// DecodeData decodes data from a io.Reader.
func DecodeData(r io.Reader) ([]any, error) {
	var data []any
	err := json.NewDecoder(r).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// NewDataFromFiles parses index data from files.
func NewDataFromFiles(d ...string) ([]any, error) {
	var data []any
	for _, datum := range d {
		p, err := dataFromFile(datum)
		if err != nil {
			return nil, err
		}
		data = append(data, p...)
	}
	return data, nil
}

func dataFromFile(d string) ([]any, error) {
	data, err := os.Open(d)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return DecodeData(data)
}

// NewDataFromString parses index data from a json formatted string.
func NewDataFromString(d string) ([]any, error) {
	buf := bytes.NewBufferString(d)
	return DecodeData(buf)
}

func parseData(d any) ([]any, error) {
	switch val := d.(type) {
	case []byte:
		return unmarshalData(val)
	case string:
		if exist(val) {
			return dataFromFile(val)
		} else {
			return unmarshalData([]byte(val))
		}
	case []any:
		return val, nil
	}
	return nil, errors.New("data couldn't be parsed")
}

func unmarshalData(d []byte) ([]any, error) {
	var data []any
	err := json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
