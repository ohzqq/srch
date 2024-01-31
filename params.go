package srch

import (
	"encoding/json"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/ohzqq/srch/txt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const (
	// search params
	Hits                 = `hits`
	AttributesToRetrieve = `attributesToRetrieve`
	Page                 = "page"
	HitsPerPage          = "hitsPerPage"
	SortFacetsBy         = `sortFacetValuesBy`
	Query                = `query`
	ParamFacets          = "facets"
	ParamFilters         = "filters"
	FacetFilters         = `facetFilters`
	ParamFullText        = `fullText`
	NbHits               = `nbHits`
	NbPages              = `nbPages`
	SortBy               = `sortBy`
	Order                = `order`

	// Settings
	SrchAttr     = `searchableAttributes`
	FacetAttr    = `attributesForFaceting`
	DataDir      = `dataDir`
	DataFile     = `dataFile`
	DefaultField = `title`

	TextAnalyzer    = "text"
	KeywordAnalyzer = "keyword"
)

type Params struct {
	Settings url.Values
	Search   url.Values
	filters  []any
}

func NewParams() *Params {
	p := &Params{
		Settings: make(url.Values),
		Search:   make(url.Values),
	}
	return p
}

func ParseParams(params any) *Params {
	q := ParseQuery(params)
	p := NewParams()
	p.Settings = GetSettings(q)
	p.Search = GetSearchParams(q)

	if p.HasFilters() {
		for _, filters := range p.Search[FacetFilters] {
			fils, err := unmarshalFilter(filters)
			if err != nil {
				break
			}
			p.filters = append(p.filters, fils...)
		}
	}
	return p
}

func GetSettings(vals url.Values) url.Values {
	settings := make(url.Values)
	for l, v := range vals {
		if l == SrchAttr ||
			l == FacetAttr ||
			l == DataDir ||
			l == DataFile {
			settings[l] = v
		}
	}
	return settings
}

func GetSearchParams(vals url.Values) url.Values {
	params := make(url.Values)
	for l, v := range vals {
		if l != SrchAttr &&
			l != FacetAttr &&
			l != DataDir &&
			l != DataFile {
			params[l] = v
		}
	}
	return params
}

func (p Params) GetSlice(key string) []string {
	if p.Settings.Has(key) {
		return p.Settings[key]
	}
	if p.Search.Has(key) {
		return p.Search[key]
	}
	return []string{}
}

func (p Params) Get(key string) string {
	if p.Settings.Has(key) {
		return p.Settings.Get(key)
	}
	if p.Search.Has(key) {
		return p.Search.Get(key)
	}
	return ""
}

func (p *Params) Set(key string, val string) {
	switch key {
	case FacetAttr, SrchAttr, DataDir, DataFile:
		p.Settings.Set(key, val)
	default:
		p.Search.Set(key, val)
	}
}

func (p *Params) AndFilter(field string, filters ...string) *Params {
	p.filters = append(p.filters, NewFilter(field, filters)...)
	return p
}

func (p *Params) OrFilter(field string, filters ...string) {
	p.filters = append(p.filters, NewFilter(field, filters))
}

func (p Params) IsFacet(attr string) bool {
	return slices.Contains(p.FacetAttr(), attr)
}

func (p *Params) NewField(attr string) *Field {
	f := NewField(attr)

	switch p.IsFacet(attr) {
	case true:
		f.SetAnalyzer(txt.Keyword())
	default:
		if p.IsFullText() {
			f.SetAnalyzer(txt.Fulltext())
		}
	}

	f.SortBy = p.SortFacetsBy()
	return f
}

func (p *Params) newFieldsMap(attrs []string) map[string]*Field {
	fields := make(map[string]*Field)
	for _, attr := range attrs {
		fields[attr] = p.NewField(attr)
	}
	return fields
}

func (p *Params) SrchAttr() []string {
	return p.GetSlice(SrchAttr)
}

func (p *Params) Fields() map[string]*Field {
	return p.newFieldsMap(p.SrchAttr())
}

func (p Params) HasFilters() bool {
	return p.Search.Has(FacetFilters)
}

func (p *Params) Facets() map[string]*Field {
	return p.newFieldsMap(p.FacetAttr())
}

func (p Params) FacetAttr() []string {
	return p.GetSlice(FacetAttr)
}

func (f Filters) String() string {
	d, err := json.Marshal(f)
	if err != nil {
		return ""
	}
	return string(d)
}

func (p *Params) SortFacetsBy() string {
	sort := SortByCount
	if p.Search.Has(SortFacetsBy) {
		if by := p.Search.Get(SortFacetsBy); by == SortByCount || by == SortByAlpha {
			sort = by
		}
	}
	return sort
}

func (p Params) Page() int {
	pn := p.Search.Get(Page)
	page, err := strconv.Atoi(pn)
	if err != nil {
		return 0
	}
	return page
}

func (p *Params) SetPage(i any) {
	p.Search.Set(Page, cast.ToString(i))
}

func (p Params) HitsPerPage() int {
	pn := p.Search.Get(HitsPerPage)
	page, err := strconv.Atoi(pn)
	if err != nil {
		return 0
	}
	return page
}

func (p Params) SetHitsPerPage(i any) {
	p.Search.Set(HitsPerPage, cast.ToString(i))
}

func (p *Params) IsFullText() bool {
	return p.Search.Has(ParamFullText)
}

func (p Params) Query() string {
	return p.Get(Query)
}

func (p Params) GetAnalyzer() string {
	if p.IsFullText() {
		return TextAnalyzer
	}
	return KeywordAnalyzer
}

func (p Params) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(p.Encode())
	if err != nil {
		return nil, err
	}
	return d, err
}

func (p *Params) UnmarshalJSON(d []byte) error {
	var pn string
	err := json.Unmarshal(d, &pn)
	if err != nil {
		return err
	}

	err = p.Decode(pn)
	if err != nil {
		return err
	}
	return nil
}

func (p Params) Encode() string {
	return p.Values().Encode()
}

func (p Params) String() string {
	return p.Values().Encode()
}

func (p Params) Values() url.Values {
	m := lo.Assign(p.Settings, p.Search)
	return url.Values(m)
}

func (p *Params) Decode(str string) error {
	q, err := url.ParseQuery(str)
	if err != nil {
		return err
	}

	p.Settings = GetSettings(q)
	p.Search = GetSearchParams(q)
	return nil
}

func ParseQuery(queries ...any) url.Values {
	q := make(url.Values)
	for _, query := range queries {
		vals, err := ParseValues(query)
		if err != nil {
			continue
		}
		for k, val := range vals {
			for _, v := range val {
				q.Add(k, v)
			}
		}
	}
	return q
}

// ParseValues takes an interface{} and returns a url.Values.
func ParseValues(f any) (url.Values, error) {
	filters := make(url.Values)
	var err error
	switch val := f.(type) {
	case url.Values:
		filters = val
	case []byte:
		filters, err = ParseQueryBytes(val)
	case string:
		filters, err = ParseQueryString(val)
	default:
		filters, err = cast.ToStringMapStringSliceE(val)
		if err != nil {
			return nil, err
		}
	}
	filters[SrchAttr] = parseSrchAttr(filters)
	filters[FacetAttr] = parseFacetAttr(filters)
	return filters, nil
}

func parseSrchAttr(vals url.Values) []string {
	if !vals.Has(SrchAttr) {
		return []string{DefaultField}
	}
	vals[SrchAttr] = GetQueryStringSlice(SrchAttr, vals)
	if len(vals[SrchAttr]) < 1 {
		vals[SrchAttr] = []string{DefaultField}
	}
	return vals[SrchAttr]
}

func parseFacetAttr(vals url.Values) []string {
	if !vals.Has(ParamFacets) {
		vals[ParamFacets] = GetQueryStringSlice(FacetAttr, vals)
	}
	return vals[ParamFacets]
}

// ParseQueryString parses an encoded filter string.
func ParseQueryString(val string) (url.Values, error) {
	q, err := url.ParseQuery(val)
	if err != nil {
		return nil, err
	}
	return q, err
}

// ParseQueryBytes parses a byte slice to url.Values.
func ParseQueryBytes(val []byte) (url.Values, error) {
	filters, err := cast.ToStringMapStringSliceE(string(val))
	if err != nil {
		return nil, err
	}
	return url.Values(filters), err
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
