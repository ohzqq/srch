package srch

import (
	"encoding/json"
	"log"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
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
	NbHits               = `nbHits`
	NbPages              = `nbPages`
	SortBy               = `sortBy`
	Order                = `order`

	// Settings
	ParamFullText = `fullText`
	SrchAttr      = `searchableAttributes`
	FacetAttr     = `attributesForFaceting`
	SortAttr      = `sortableAttributes`
	DataDir       = `dataDir`
	DataFile      = `dataFile`
	IndexPath     = `indexPath`
	DefaultField  = `title`
	UID           = `uid`

	TextAnalyzer    = "text"
	KeywordAnalyzer = "keyword"
)

var paramsSettings = []string{
	SrchAttr,
	FacetAttr,
	SortAttr,
	DataDir,
	DataFile,
	IndexPath,
	DefaultField,
	ParamFullText,
	UID,
}

var paramsSearch = []string{
	Hits,
	AttributesToRetrieve,
	Page,
	HitsPerPage,
	SortFacetsBy,
	Query,
	ParamFacets,
	ParamFilters,
	FacetFilters,
	NbHits,
	NbPages,
	SortBy,
	Order,
}

type Params struct {
	Settings     url.Values
	Search       url.Values
	filters      []any
	FacetAttrs   []string `mapstructure:"attributesForFaceting"`
	SrchAttrs    []string `mapstructure:"searchableAttributes"`
	SortAttrs    []string `mapstructure:"sortableAttributes"`
	QueryStr     []string `mapstructure:"query"`
	ParamFacets  []string `mapstructure:"facets"`
	ParamFilters []string `mapstructure:"filters"`
	FacetFilters []string `mapstructure:"facetFilters"`
	DataDir      []string `mapstructure:"dataDir"`
	DataFile     []string `mapstructure:"dataFile"`
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

	return p
}

func GetSettings(vals url.Values) url.Values {
	settings := lo.PickByKeys(vals, paramsSettings)
	return url.Values(settings)
}

func GetSearchParams(vals url.Values) url.Values {
	params := lo.PickByKeys(vals, paramsSearch)
	return params
}

func (p Params) GetSlice(key string) []string {
	if p.Settings.Has(key) {
		return p.Settings[key]
	}
	if p.Search.Has(key) {
		return p.Search[key]
	}
	return nil
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
	case FacetAttr, SrchAttr, DataDir, DataFile, SortAttr, IndexPath, ParamFullText:
		p.Settings.Set(key, val)
	default:
		p.Search.Set(key, val)
	}
}

func (p Params) GetUID() string {
	return p.Settings.Get("uid")
}

func (p Params) Has(key string) bool {
	switch key {
	case FacetAttr, SrchAttr, DataDir, DataFile, SortAttr, IndexPath, ParamFullText:
		return p.Settings.Has(key)
	default:
		return p.Search.Has(key)
	}
}

func (p *Params) SetSearch(params string) *Params {
	q := ParseQuery(params)
	p.Search = q
	return p
}

func (p Params) HasFilters() bool {
	return p.Search.Has(FacetFilters)
}

func (p *Params) Filters() []any {
	if p.HasFilters() {
		fils, err := unmarshalFilter(p.Get(FacetFilters))
		if err != nil {
		}
		return fils
	}
	return p.filters
}

func (p *Params) SetFilters(filters []any) *Params {
	if len(filters) == 0 {
		p.Search.Del(FacetFilters)
		return p
	}

	d, err := json.Marshal(filters)
	if err != nil {
		return p
	}
	p.Set(FacetFilters, string(d))

	return p
}

func (p Params) IsFacet(attr string) bool {
	return slices.Contains(p.FacetAttr(), attr)
}

func (p *Params) NewField(attr string) *Field {
	f := NewField(attr)

	switch p.IsFacet(attr) {
	case true:
	default:
	}

	f.SortBy = p.SortFacetsBy()
	return f
}

func (p *Params) Fields() map[string]*Field {
	return p.newFieldsMap(p.SrchAttr())
}

func (p *Params) Facets() map[string]*Field {
	return p.newFieldsMap(p.FacetAttr())
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

func (p Params) FacetAttr() []string {
	return p.GetSlice(FacetAttr)
}

func (p Params) SortAttr() []string {
	sort := p.GetSlice(SortAttr)
	if sort != nil {
		return sort
	}
	return []string{"title:text"}
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
	page := viper.GetInt(HitsPerPage)
	if p.Has(HitsPerPage) {
		pn := p.Search.Get(HitsPerPage)
		page, err := strconv.Atoi(pn)
		if err != nil {
			return 25
		}
		return page
	}
	return page
}

func (p Params) SetHitsPerPage(i any) {
	p.Search.Set(HitsPerPage, cast.ToString(i))
}

func (p *Params) IsFullText() bool {
	return p.Has(ParamFullText)
}

func (p *Params) GetFullText() string {
	return p.Get(ParamFullText)
}

func (p Params) Query() string {
	return p.Get(Query)
}

func (p Params) SortBy() string {
	if p.Has(SortBy) {
		return p.Get(SortBy)
	}
	return DefaultField
}

func (p Params) GetAnalyzer() string {
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
	return lo.Assign(p.Settings, p.Search)
}

func (p *Params) HasData() bool {
	return p.Has(DataFile) ||
		p.Has(DataDir)
}

func (p Params) GetData() string {
	switch {
	case p.Has(DataFile):
		return p.Get(DataFile)
	case p.Has(DataDir):
		return p.Get(DataDir)
	default:
		return ""
	}
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

func parseSearchParamsJSON(sp any) url.Values {
	q := make(url.Values)

	var fs []byte
	switch val := sp.(type) {
	case string:
		fs = []byte(val)
	case []byte:
		fs = val
	}

	raw := make(map[string]any)
	err := json.Unmarshal(fs, &raw)
	if err != nil {
		log.Fatal(err)
		return q
	}

	for key, param := range raw {
		switch val := param.(type) {
		case []any:
			if key == FacetFilters {
				d, err := json.Marshal(val)
				if err != nil {
					break
				}
				q.Set(FacetFilters, string(d))
				break
			}
			for _, v := range val {
				q.Add(key, cast.ToString(v))
			}
		case float64:
			q.Set(key, cast.ToString(val))
		case bool:
			q.Set(key, "")
		case string:
			q.Set(key, val)
		}
	}

	return q
}

func ParseSearchParamsJSON(sp any) string {
	return parseSearchParamsJSON(sp).Encode()
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
