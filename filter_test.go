package srch

import (
	"log"
	"net/url"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/ohzqq/audible"
	"github.com/samber/lo"
)

func TestAudibleSearch(t *testing.T) {
	cfg := map[string]any{
		"searchableFields": []string{"Title"},
	}
	a, err := New(cfg, WithSearch(audibleSrch), Interactive)
	if err != nil {
		t.Error(err)
	}
	v := make(url.Values)
	v.Set("q", "amy lane fish")
	res := a.Search(v)
	println(res.Len())
	//res.Print()
}

func audibleSrch(q string) []any {
	s := audible.NewSearch(q)
	r, err := audible.Products().Search(s)
	if err != nil {
		log.Fatal(err)
	}
	var sl []map[string]any
	for _, p := range r.Products {
		a := make(map[string]any)
		mapstructure.Decode(p, &a)
		sl = append(sl, a)
	}
	//fmt.Printf("products %v\n", r.Products)
	return lo.ToAnySlice(sl)
}
