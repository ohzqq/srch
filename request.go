package srch

import (
	"net/url"
	"strings"

	"github.com/ohzqq/srch/param"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Request struct {
	*param.Params
	filters []string
}

func NewRequest() *Request {
	return &Request{
		Params: param.New(),
	}
}

func ParseRequest(params string) (*Request, error) {
	p, err := param.Parse(params)
	if err != nil {
		return nil, err
	}
	return &Request{Params: p}, nil
}

func (r *Request) Parse(params string) (*Request, error) {
	return ParseRequest(params)
}

func (r *Request) SetValues(vals url.Values) *Request {
	r.Params.Set(vals)
	return r
}

func (p *Request) SetRoute(path string) *Request {
	p.URL.Path = path
	p.Params.Route = strings.TrimPrefix(path, "/")
	return p
}

func (p *Request) SetPath(path string) *Request {
	p.Params.Path = path
	return p
}

func (r *Request) SrchAttr(attr ...string) *Request {
	r.Params.SrchAttr = attr
	return r
}

func (r *Request) FacetAttr(attr ...string) *Request {
	r.Params.FacetAttr = attr
	return r
}

func (r *Request) SortAttr(attr ...string) *Request {
	r.Params.SortAttr = attr
	return r
}

func (r *Request) RtrvAttr(attr ...string) *Request {
	//r.Params.RtrvAttr = attr
	return r
}

func (r *Request) Facets(attr ...string) *Request {
	r.Params.Facets = attr
	return r
}

func (r *Request) OrFilter(filters ...string) *Request {
	return r.FacetFilters(lo.ToAnySlice(filters))
}

func (r *Request) AndFilter(filters ...string) *Request {
	return r.FacetFilters(lo.ToAnySlice(filters)...)
}

func (r *Request) FacetFilters(filters ...any) *Request {
	r.Params.FacetFilters = append(r.Params.FacetFilters, filters...)
	return r
}

func (r *Request) Query(val string) *Request {
	r.Params.Query = val
	return r
}

func (r *Request) Filters(val string) *Request {
	r.Params.Filters = val
	return r
}

func (r *Request) SortFacetsBy(val string) *Request {
	r.Params.SortFacetsBy = val
	return r
}

func (r *Request) SortBy(val string) *Request {
	r.Params.SortBy = val
	return r
}

func (r *Request) Order(val string) *Request {
	r.Params.Order = val
	return r
}

func (r *Request) Format(val string) *Request {
	r.Params.Format = val
	return r
}

func (r *Request) DefaultField(val string) *Request {
	r.Params.DefaultField = val
	return r
}

func (r *Request) UID(val string) *Request {
	r.Params.UID = val
	return r
}

func (r *Request) Page(p int) *Request {
	r.Params.Page = p
	return r
}

func (r *Request) HitsPerPage(p int) *Request {
	r.Params.HitsPerPage = p
	return r
}

func GetViperParams() *Request {
	vals := viper.AllSettings()

	params := make(url.Values)
	for key, val := range cast.ToStringMapStringSlice(vals) {
		for _, k := range param.SettingParams {
			if key == strings.ToLower(k) {
				params[k] = val
			}
		}
		for _, k := range param.SearchParams {
			if key == strings.ToLower(k) {
				params[k] = val
			}
		}
	}

	req := NewRequest().SetValues(params)

	for _, key := range param.Routes {
		if viper.IsSet(strings.ToLower(key)) {
			val := viper.GetStringSlice(key)
			req.SetRoute(key)
			req.SetPath(val[0])
		}
	}

	return req
}

func init() {
}
