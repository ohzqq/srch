package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
)

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
	mime.AddExtensionType(".bleve", "application/bleve")
}

const (
	NdJSON = `application/x-ndjson`
	JSON   = `application/json`
)

type Doc interface {
	GetID() int
}

type Source struct {
	docs []Doc
}

type Data struct {
	Path  string
	Route string
	Files []*File
	data  []map[string]any
	ids   []string
}

type File struct {
	Path        string
	ContentType string
}

func New(r, path string) *Data {
	d := &Data{
		Route: r,
		Path:  path,
	}

	switch route := d.Route; route {
	case "dir":
		js, _ := filepath.Glob(filepath.Join(d.Path, "*.json"))
		d.AddFile(js...)

		nd, _ := filepath.Glob(filepath.Join(d.Path, "*.ndjson"))
		d.AddFile(nd...)
	default:
		d.AddFile(d.Path)
	}

	return d
}

func NewData() *Data {
	return &Data{}
}

func (d *Data) AddFile(paths ...string) {
	for _, path := range paths {
		file, err := NewFile(path)
		if err != nil {
			return
		}
		d.Files = append(d.Files, file)
	}
}

func (d *Data) decodeNDJSON() error {
	i := 0
	for _, file := range d.Files {
		f, err := os.Open(file.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		dec := json.NewDecoder(f)
		for {
			m := make(map[string]any)
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			var id any
			if uid, ok := m["id"]; ok {
				id = uid
			} else {
				id = i
			}
			d.ids = append(d.ids, cast.ToString(id))
			d.data = append(d.data, m)
			i++
		}
		//return nil
	}
	return nil
}

func (d *Data) Decode() ([]map[string]any, error) {
	for _, file := range d.Files {
		f, err := os.Open(file.Path)
		if err != nil {
			return d.data, err
		}
		defer f.Close()

		switch file.ContentType {
		case NdJSON:
			err := DecodeNDJSON(f, &d.data)
			if err != nil {
				return nil, err
			}
		case JSON:
			err := DecodeJSON(f, &d.data)
			if err != nil {
				return nil, err
			}
		}
	}

	return d.data, nil
}

func NewFile(path string) (*File, error) {
	file := &File{
		Path: path,
	}

	ext := filepath.Ext(file.Path)
	if ext == "" {
		return file, errors.New("file has no extension")
	}

	file.ContentType = mime.TypeByExtension(ext)
	if b, _, ok := strings.Cut(file.ContentType, ";"); ok {
		file.ContentType = b
	}

	return file, nil
}

func DecodeData(r io.Reader, ct string, data *[]map[string]any) error {
	switch ct {
	case "application/x-ndjson":
		return DecodeNDJSON(r, data)
	case "application/json":
		return DecodeJSON(r, data)
	}
	return nil
}

// DecodeNDJSON decodes data from a io.Reader.
func DecodeNDJSON(r io.Reader, data *[]map[string]any) error {
	dec := json.NewDecoder(r)
	for {
		m := make(map[string]any)
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		*data = append(*data, m)
	}
	return nil
}

func DecodeJSON(r io.Reader, data *[]map[string]any) error {
	dec := json.NewDecoder(r)
	for {
		m := []map[string]any{}
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		*data = append(*data, m...)
	}
	return nil
}

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
	return DecodeDataO(r)
}

// DirSrc parses json files from a directory.
func DirSrc(dir string) ([]map[string]any, error) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	files, err := filepath.Glob(dir + "*.json")
	if err != nil {
		return nil, err
	}
	return FileSrc(files...)
}

// StringSrc parses raw json formatted data.
func ByteSrc(d []byte) ([]map[string]any, error) {
	return DecodeDataO(bytes.NewBuffer(d))
}

// StringSrc parses index data from a json formatted string.
func StringSrc(d string) ([]map[string]any, error) {
	return DecodeDataO(bytes.NewBufferString(d))
}

// DecodeDataO decodes data from a io.Reader.
func DecodeDataO(r io.Reader) ([]map[string]any, error) {
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
	return DecodeDataO(data)
}
