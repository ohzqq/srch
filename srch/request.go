package srch

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ohzqq/srch/endpoint"
	"github.com/samber/lo"
)

type Request struct {
	*url.URL

	vals      url.Values
	wildCards []string
	method    string
}

func NewRequest(u any) (*Request, error) {
	v, err := ParseQuery(u)
	if err != nil {
		return nil, err
	}
	return &Request{vals: v}, nil
}

func ParseHTTPRequest(r *http.Request) *Request {
	r.ParseForm()
	cards := getWildCards(r)

	params := lo.Assign(
		map[string][]string(r.Form),
		map[string][]string(r.PostForm),
		map[string][]string(r.URL.Query()),
	)

	return &Request{
		vals:      params,
		URL:       r.URL,
		method:    r.Method,
		wildCards: cards,
	}
}

func (req *Request) Endpoint() endpoint.Endpoint {
	return endpoint.Parse(req.Path, req.wildCards)
}

func (req *Request) Route() string {
	return req.Endpoint().Route()
}

func (req *Request) Cfg() (*Cfg, error) {
	return NewCfg(req.vals)
}

func (req *Request) Client() (*Client, error) {
	cfg, err := req.Cfg()
	if err != nil {
		return nil, err
	}
	return NewClient(cfg)
}

func Indexes(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func IdxSrch(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func IdxBrowse(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func IdxObject(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func IdxQuery(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func IdxSettings(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func Facets(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func Facet(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func FacetQuery(w http.ResponseWriter, r *http.Request) {
	req := ParseHTTPRequest(r)
	fmt.Fprintf(w, "%#v", req)
}

func ParseQuery(q any) (url.Values, error) {
	switch v := q.(type) {
	case string:
		v = strings.TrimPrefix(v, "?")
		return url.ParseQuery(v)
	case map[string][]string:
		return url.Values(v), nil
	case url.Values:
		return v, nil
	case *url.URL:
		return v.Query(), nil
	default:
		return nil, errors.New("param must be of type: string, map[string][]string, url.Values, *url.URL")
	}
}
