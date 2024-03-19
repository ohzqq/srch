package param

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func init() {
	mime.AddExtensionType(".ndjson", "application/x-ndjson")
}

type Params struct {
	URL   *url.URL   `json:"-"`
	Other url.Values `json:"-"`

	// Index Settings
	SrchAttr     []string `query:"searchableAttributes,omitempty" json:"searchableAttributes,omitempty"`
	FacetAttr    []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty"`
	SortAttr     []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty"`
	DefaultField string   `query:"defaultField,omitempty" json:"defaultField,omitempty"`
	UID          string   `query:"uid,omitempty" json:"uid,omitempty"`

	// Search
	Hits                 int      `query:"hits,omitempty" json:"hits,omitempty"`
	AttributesToRetrieve []string `query:"attributesToRetrieve,omitempty" json:"attributesToRetrieve,omitempty"`
	Page                 int      `query:"page,omitempty" json:"page,omitempty"`
	HitsPerPage          int      `query:"hitsPerPage,omitempty" json:"hitsPerPage,omitempty"`
	Query                string   `query:"query,omitempty" json:"query,omitempty"`
	SortBy               string   `query:"sortBy,omitempty" json:"sortBy,omitempty"`
	Order                string   `query:"order,omitempty" json:"order,omitempty"`
	// Facets
	Facets       []string `query:"facets,omitempty" json:"facets,omitempty"`
	Filters      string   `query:"filters,omitempty" json:"filters,omitempty"`
	FacetFilters []any    `query:"facetFilters,omitempty" json:"facetFilters,omitempty"`
	SortFacetsBy string   `query:"sortFacetsBy,omitempty" json:"sortFacetsBy,omitempty"`

	// Data
	Format string `json:"-"`
	Path   string `json:"-"`
	Route  string `json:"-"`
}

var (
	blvPath  = regexp.MustCompile(`/?(blv)/(.*)`)
	filePath = regexp.MustCompile(`/?(file)/(.*)`)
	dirPath  = regexp.MustCompile(`/?(dir)/(.*)`)
	routes   = []*regexp.Regexp{blvPath, filePath, dirPath}
)

func New() *Params {
	return &Params{
		URL:   &url.URL{},
		Other: make(url.Values),
	}
}

func Parse(params string) (*Params, error) {
	p := New()

	u, err := url.Parse(params)
	if err != nil {
		return nil, fmt.Errorf("param url parsing err: %w\n", err)
	}
	p.URL = u

	p.ParseRoute(p.URL.Path)

	err = p.Set(u.Query())
	if err != nil {
		return nil, fmt.Errorf("error setting params %w\n", err)
	}

	return p, nil
}

func (p Params) Search() url.Values {
	vals := make(url.Values)
	for _, key := range paramsSearch {
		q := p.URL.Query()
		if q.Has(key) {
			vals[key] = q[key]
		}
	}
	return vals
}

func (p Params) Index() url.Values {
	vals := make(url.Values)
	for _, key := range paramsSettings {
		q := p.URL.Query()
		if q.Has(key) {
			vals[key] = q[key]
		}
	}
	return vals
}

func (p *Params) ParseRoute(route string) *Params {
	p.Route, p.Path = parseRoute(route)
	return p
}

func (p *Params) parseRoute() *Params {
	p.Route, p.Path = parseRoute(p.URL.Path)
	return p
}

func parseRoute(path string) (string, string) {
	for _, reg := range routes {
		matches := reg.FindStringSubmatch(path)
		if len(matches) > 1 {
			pre := matches[1]
			loc := matches[2]
			loc, err := filepath.Abs(loc)
			if err != nil {
				loc = ""
			}
			return pre, loc
		}
	}
	return "", ""
}

func (s *Params) Set(v url.Values) error {
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
		case UID:
			s.UID = v.Get(key)
		}
	}
	for _, key := range paramsData {
		switch key {
		case Format:
			if v.Has(key) {
				s.Format = v.Get(key)
			}
		}
	}
	for _, key := range paramsSearch {
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
					return fmt.Errorf("failed to unmarshal filters %v\nerr: %w\n", fil, err)
				}
				s.FacetFilters = f
			}
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
	case Format:
		return s.Format != ""
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
		case UID:
			vals.Set(key, s.UID)
		}
	}
	for _, key := range paramsData {
		if !s.Has(key) {
			continue
		}
		switch key {
		case Format:
			vals.Set(key, s.Format)
		}
	}
	for _, key := range paramsSearch {
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
			d, err := json.Marshal(s.FacetFilters)
			if err == nil {
				vals.Set(key, string(d))
			}
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

func (p *Params) Encode() string {
	return p.Values().Encode()
}

func (p *Params) String() string {
	p.URL.RawQuery = p.Encode()
	return p.URL.String()
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
