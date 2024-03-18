package param

import (
	"encoding/json"
	"mime"
	"net/url"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
}

type Params struct {
	// Settings

	// Index Settings
	SrchAttr     []string `query:"searchableAttributes,omitempty" json:"searchableAttributes,omitempty"`
	FacetAttr    []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty"`
	SortAttr     []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty"`
	DefaultField string   `query:"defaultField,omitempty" json:"defaultField,omitempty"`

	// Search
	Hits                 int      `query:"hits,omitempty" json:"hits,omitempty"`
	AttributesToRetrieve []string `query:"attributesToRetrieve,omitempty" json:"attributesToRetrieve,omitempty"`
	Page                 int      `query:"page,omitempty" json:"page,omitempty"`
	HitsPerPage          int      `query:"hitsPerPage,omitempty" json:"hitsPerPage,omitempty"`
	Query                string   `query:"query,omitempty" json:"query,omitempty"`
	SortBy               string   `query:"sortBy,omitempty" json:"sortBy,omitempty"`
	Order                string   `query:"order,omitempty" json:"order,omitempty"`

	// Cfg
	BlvPath  string `query:"fullText,omitempty" json:"fullText,omitempty"`
	DataDir  string `query:"dataDir,omitempty" json:"dataDir,omitempty"`
	DataFile string `query:"dataFile,omitempty" json:"dataFile,omitempty"`
	UID      string `query:"uid,omitempty" json:"uid,omitempty"`

	// Facets
	Facets       []string `query:"facets,omitempty" json:"facets,omitempty"`
	Filters      string   `query:"filters,omitempty" json:"filters,omitempty"`
	FacetFilters []any    `query:"facetFilters,omitempty" json:"facetFilters,omitempty"`
	SortFacetsBy string   `query:"sortFacetsBy,omitempty" json:"sortFacetsBy,omitempty"`
	params       url.Values
	Other        url.Values `mapstructure:"params,omitempty"`
}

func New() *Params {
	return &Params{
		Other:  make(url.Values),
		params: make(url.Values),
	}
}

func Parse(params string) (*Params, error) {
	p := New()

	if !strings.HasPrefix(params, "?") {
		params = "?" + params
	}

	vals, err := url.Parse(params)
	if err != nil {
		return nil, err
	}

	err = p.Set(vals.Query())
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Params) IsFile() bool {
	return p.Has(DataDir) ||
		p.Has(DataFile) ||
		p.Has(BlvPath)
}

func (s *Params) Set(v url.Values) error {
	s.params = v
	for _, key := range paramsSettings {
		switch key {
		case SrchAttr:
			s.SrchAttr = parseSrchAttr(v)
		case FacetAttr:
			s.FacetAttr = parseFacetAttr(v)
		case SortAttr:
			s.SortAttr = GetQueryStringSlice(key, v)
		case DefaultField:
			s.DefaultField = v.Get(key)
		}
	}
	for _, key := range paramsCfg {
		switch key {
		case DataDir:
			s.DataDir = v.Get(key)
		case DataFile:
			s.DataFile = v.Get(key)
		case FullText:
			s.BlvPath = v.Get(key)
		case UID:
			s.UID = v.Get(key)
		}
	}
	for _, key := range paramsFacets {
		switch key {
		case SortFacetsBy:
			s.SortFacetsBy = v.Get(key)
		case Facets:
			s.Facets = GetQueryStringSlice(key, v)
		case Filters:
			s.Filters = v.Get(key)
		case FacetFilters:
			if v.Has(key) {
				fil := v.Get(key)
				f, err := unmarshalFilter(fil)
				if err != nil {
					return err
				}
				s.FacetFilters = f
			}
		}
	}
	for _, key := range paramsSearch {
		switch key {
		case Hits:
			s.Hits = GetQueryInt(key, v)
		case AttributesToRetrieve:
			s.AttributesToRetrieve = GetQueryStringSlice(key, v)
		case Page:
			s.Page = GetQueryInt(key, v)
		case HitsPerPage:
			s.HitsPerPage = GetQueryInt(key, v)
		case Query:
			s.Query = v.Get(key)
		case SortBy:
			s.SortBy = v.Get(key)
		case Order:
			s.Order = v.Get(key)
		}
	}
	return nil
}

func (s *Params) Has(key string) bool {
	switch key {
	case Hits:
		return s.Hits != 0
	case AttributesToRetrieve:
		return len(s.AttributesToRetrieve) != 0
	case Page:
		return s.Page != 0
	case HitsPerPage:
		return s.HitsPerPage != 0
	case Query:
		return s.Query != ""
	case SortBy:
		return s.SortBy != ""
	case Order:
		return s.Order != ""
	case DataDir:
		return s.DataDir != ""
	case DataFile:
		return s.DataFile != ""
	case FullText:
		return s.BlvPath != ""
	case UID:
		return s.UID != ""
	case SortFacetsBy:
		return s.SortFacetsBy != ""
	case Facets:
		return len(s.Facets) > 0
	case Filters:
		return s.Filters != ""
	case FacetFilters:
		return len(s.FacetFilters) > 0
	case SrchAttr:
		return len(s.SrchAttr) > 0
	case FacetAttr:
		return len(s.FacetAttr) > 0
	case SortAttr:
		return len(s.SortAttr) > 0
	case DefaultField:
		return s.DefaultField != ""
	default:
		return false
	}
}

func (s *Params) Values() url.Values {
	vals := make(url.Values)
	for _, key := range paramsSettings {
		if !s.Has(key) {
			continue
		}
		switch key {
		case SrchAttr:
			vals[key] = s.SrchAttr
		case FacetAttr:
			vals[key] = s.FacetAttr
		case SortAttr:
			vals[key] = s.SortAttr
		case DefaultField:
			vals.Set(key, s.DefaultField)
		}
	}
	for _, key := range paramsCfg {
		if !s.Has(key) {
			continue
		}
		switch key {
		case DataDir:
			vals.Set(key, s.DataDir)
		case DataFile:
			vals.Set(key, s.DataFile)
		case FullText:
			vals.Set(key, s.BlvPath)
		case UID:
			vals.Set(key, s.UID)
		}
	}
	for _, key := range paramsFacets {
		if !s.Has(key) {
			continue
		}
		switch key {
		case SortFacetsBy:
			vals.Set(key, s.SortFacetsBy)
		case Facets:
			vals[key] = s.Facets
		case Filters:
			vals.Set(key, s.Filters)
		case FacetFilters:
			for _, f := range s.FacetFilters {
				vals.Add(key, cast.ToString(f))
			}
		}
	}
	for _, key := range paramsSearch {
		if !s.Has(key) {
			continue
		}
		switch key {
		case Hits:
			vals.Set(key, cast.ToString(s.Hits))
		case AttributesToRetrieve:
			vals[key] = s.AttributesToRetrieve
		case Page:
			vals.Set(key, cast.ToString(s.Page))
		case HitsPerPage:
			vals.Set(key, cast.ToString(s.HitsPerPage))
		case Query:
			vals.Set(key, s.Query)
		case SortBy:
			vals.Set(key, s.SortBy)
		case Order:
			vals.Set(key, s.Order)
		}
	}
	return vals
}

func (p Params) IsFullText() bool {
	return p.BlvPath != ""
}

func (p Params) HasData() bool {
	return p.Has(DataDir) || p.Has(DataFile)
}

func (p Params) GetDataFiles() []string {
	var data []string
	switch {
	case p.Has(DataDir):
		data = append(data, p.DataDir)
	case p.Has(DataFile):
		data = append(data, p.DataFile)
	}
	return data
}

func (p *Params) Encode() string {
	return p.Values().Encode()
}

func (p *Params) String() string {
	return p.Values().Encode()
}

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (url.Values, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, err
}

func GetAnySlice(key string, vals url.Values) []any {
	return lo.ToAnySlice(GetQueryStringSlice(key, vals))
}

func GetQueryInt(key string, vals url.Values) int {
	if vals.Has(key) {
		return cast.ToInt(vals.Get(key))
	}
	return 0
}

func GetQueryStringSlice(key string, q url.Values) []string {
	var vals []string
	if q.Has(key) {
		for _, val := range q[key] {
			if val == "" {
				break
			}
			for _, v := range strings.Split(val, ",") {
				vals = append(vals, v)
			}
		}
	}
	return vals
}

func unmarshalFilter(dec string) ([]any, error) {
	var f []any
	err := json.Unmarshal([]byte(dec), &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func parseSrchAttr(vals url.Values) []string {
	if !vals.Has(SrchAttr) {
		return []string{"*"}
	}
	v := GetQueryStringSlice(SrchAttr, vals)
	if len(v) > 0 {
		return v
	}
	return []string{"*"}
}

func parseFacetAttr(vals url.Values) []string {
	if !vals.Has(FacetAttr) {
		vals[FacetAttr] = GetQueryStringSlice(FacetAttr, vals)
	}
	return vals[Facets]
}
