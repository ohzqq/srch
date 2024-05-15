package param

import (
	"net/url"
	"slices"
	"testing"

	"github.com/ohzqq/sp"
)

const urlq = `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&index=default`

var testParams = &Params{
	SrchAttr:  []string{"title"},
	FacetAttr: []string{"tags", "authors", "series", "narrators"},
	Index:     "default",
}

func TestUnmarshal(t *testing.T) {
	q := parsed()
	p := Params{}
	err := sp.Decode(q, &p)
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

func TestMarshal(t *testing.T) {
	v, err := sp.Encode(testParams)
	if err != nil {
		t.Error(err)
	}
	//fmt.Printf("%#v\n", v)
	sw := []string{"title"}
	if !slices.Equal(v["searchableAttributes"], sw) {
		t.Errorf("got %v, expected %v\n", v["searchableAttributes"], sw)
	}
	facets := []string{"tags", "authors", "series", "narrators"}
	if !slices.Equal(v["attributesForFaceting"], facets) {
		t.Errorf("got %v, expected %v\n", v["attributesForFaceting"], facets)
	}
	i := []string{"default"}
	if !slices.Equal(v["index"], i) {
		t.Errorf("got %v, expected %v\n", v["index"], i)
	}
}

func parsed() url.Values {
	v, _ := url.ParseQuery(urlq)
	//fmt.Printf("%#v\n", v)
	return v
}
