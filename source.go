package srch

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// StringSliceSrc takes a string slice and returns data for for indexing with
// the default field of 'title'.
func StringSliceSrc(data []string) []map[string]any {
	d := make([]map[string]any, len(data))
	for i, item := range data {
		d[i] = map[string]any{"title": item}
	}
	return d
}

// FileSrc takes json data files.
func FileSrc(files ...string) ([]map[string]any, error) {
	var data []map[string]any
	for _, file := range files {
		p, err := dataFromFile(file)
		if err != nil {
			return nil, err
		}
		data = append(data, p...)
	}
	return data, nil
}

// ReaderSrc takes an io.Reader of a json stream.
func ReaderSrc(r io.Reader) ([]map[string]any, error) {
	return DecodeData(r)
}

// DirSrc parses json files from a directory.
func DirSrc(dir string) ([]map[string]any, error) {
	files, err := filepath.Glob(dir + "*.json")
	if err != nil {
		return nil, err
	}
	return FileSrc(files...)
}

// StringSrc parses raw json formatted data.
func ByteSrc(d []byte) ([]map[string]any, error) {
	return DecodeData(bytes.NewBuffer(d))
}

// StringSrc parses index data from a json formatted string.
func StringSrc(d string) ([]map[string]any, error) {
	return DecodeData(bytes.NewBufferString(d))
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

func dataFromFile(d string) ([]map[string]any, error) {
	data, err := os.Open(d)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return DecodeData(data)
}
