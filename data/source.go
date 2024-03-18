package data

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ohzqq/srch/param"
)

const (
	NdJSON = `application/x-ndjson`
	JSON   = `application/json`
)

func Get(data *[]map[string]any, paths ...string) error {
	wfn := func(path string, info fs.DirEntry, err error) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		name := info.Name()
		ct := mime.TypeByExtension(filepath.Ext(name))
		if b, _, ok := strings.Cut(ct, ";"); ok {
			ct = b
		}

		switch ct {
		case "application/x-ndjson":
			return DecodeNDJSON(f, data)
		case "application/json":
			return DecodeJSON(f, data)
		}
		return nil
	}
	for _, p := range paths {
		u, err := url.Parse(p)
		if err != nil {
			break
		}
		if s := u.Scheme; s == "" || s == "file" {
			p := strings.TrimPrefix(p, "file://")
			return filepath.WalkDir(p, wfn)
		}
	}
	return nil
}

func FileTransport(path string) http.RoundTripper {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir(path)))
	return t
}

func NewClient(params string) (*http.Client, error) {
	p, err := param.Parse(params)
	if err != nil {
		return nil, err
	}

	var path string
	switch {
	case p.Has(param.DataDir):
		path = p.SrchCfg.DataDir
	case p.Has(param.DataFile):
		path = p.SrchCfg.DataFile[0]
	case p.Has(param.BlvPath):
		path = p.SrchCfg.BlvPath
	}

	req := &http.Client{}

	if path != "" {
		req.Transport = FileTransport(path)
	}

	return req, nil
}

func GetFSData(data *[]map[string]any, paths ...string) error {
	wfn := func(path string, info fs.DirEntry, err error) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		name := info.Name()
		ct := mime.TypeByExtension(filepath.Ext(name))
		if b, _, ok := strings.Cut(ct, ";"); ok {
			ct = b
		}

		switch ct {
		case "application/x-ndjson":
			return DecodeNDJSON(f, data)
		case "application/json":
			return DecodeJSON(f, data)
		}
		return nil
	}
	for _, p := range paths {
		return filepath.WalkDir(p, wfn)
	}
	return nil
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
