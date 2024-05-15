package param

import (
	"net/url"
	"slices"
	"testing"

	"github.com/ohzqq/srch/param"
	"github.com/sonh/qs"
)

const urlq = `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&index=default`

var testParams = &param.Params{
	SrchAttr:  []string{"title"},
	FacetAttr: []string{"tags", "authors", "series", "narrators"},
	Index:     "default",
}

func TestBind(t *testing.T) {
	q := parsed()
	b := DefaultBinder{}
	p := param.Params{}
	err := b.BindQueryParams(q, &p)
	if err != nil {
		t.Error(err)
	}
	sw := []string{"title"}
	if !slices.Equal(p.SrchAttr, sw) {
		t.Errorf("got %v, expected %v\n", p.SrchAttr, sw)
	}
	facets := []string{"tags,authors,series,narrators"}
	if !slices.Equal(p.FacetAttr, facets) {
		t.Errorf("got %v, expected %v\n", p.FacetAttr, facets)
	}
	i := "default"
	if p.Index != i {
		t.Errorf("got %v, expected %v\n", p.Index, i)
	}
}

func TestUnmarshal(t *testing.T) {
	enc := qs.NewEncoder()
	v, err := enc.Values(testParams)
	if err != nil {
		t.Error(err)
	}
	//fmt.Printf("%#v\n", v)
	sw := []string{"title"}
	if !slices.Equal(v[param.SrchAttr.String()], sw) {
		t.Errorf("got %v, expected %v\n", v[param.SrchAttr.String()], sw)
	}
	facets := []string{"tags", "authors", "series", "narrators"}
	if !slices.Equal(v[param.FacetAttr.String()], facets) {
		t.Errorf("got %v, expected %v\n", v[param.FacetAttr.String()], facets)
	}
	i := []string{"default"}
	if !slices.Equal(v[param.Index.String()], i) {
		t.Errorf("got %v, expected %v\n", v[param.Index.String()], i)
	}
}

func parsed() url.Values {
	v, _ := url.ParseQuery(urlq)
	//fmt.Printf("%#v\n", v)
	return v
}
