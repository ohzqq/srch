package param

import (
	"net/url"
	"slices"
	"testing"
)

const urlq = `searchableAttributes=title&attributesForFaceting=tags,authors,series,narrators&index=default`

var testParams = &Cfg{
	SrchAttr:  []string{"title"},
	FacetAttr: []string{"tags", "authors", "series", "narrators"},
	Index:     "default",
}

func TestDecode(t *testing.T) {
	cfg := &Cfg{}
	err := Decode("?"+urlq, cfg)
	if err != nil {
		t.Error(err)
	}
	sw := []string{"title"}
	if !slices.Equal(cfg.SrchAttr, sw) {
		t.Errorf("got %v, expected %v\n", cfg.SrchAttr, sw)
	}
	facets := []string{"tags", "authors", "series", "narrators"}
	if !slices.Equal(cfg.FacetAttr, facets) {
		t.Errorf("got %v, expected %v\n", cfg.FacetAttr, facets)
	}
	i := "default"
	if cfg.Index != i {
		t.Errorf("got %v, expected %v\n", cfg.Index, i)
	}
}

func TestDecodeCfg(t *testing.T) {
	cfg, err := ParseCfg("?" + urlq)
	if err != nil {
		t.Error(err)
	}
	sw := []string{"title"}
	if !slices.Equal(cfg.SrchAttr, sw) {
		t.Errorf("got %v, expected %v\n", cfg.SrchAttr, sw)
	}
	facets := []string{"tags", "authors", "series", "narrators"}
	if !slices.Equal(cfg.FacetAttr, facets) {
		t.Errorf("got %v, expected %v\n", cfg.FacetAttr, facets)
	}
	i := "default"
	if cfg.Index != i {
		t.Errorf("got %v, expected %v\n", cfg.Index, i)
	}
}

func TestEncodeCfg(t *testing.T) {
	v, err := testParams.Encode()
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
