package srch

import (
	"log"
	"net/url"
	"testing"

	"github.com/ohzqq/audible"
	"github.com/samber/lo"
)

func TestAudibleSearch(t *testing.T) {
	cfg := map[string]any{
		"searchableFields": []string{"title"},
	}
	a, err := New(cfg, WithSearch(audibleSrch), Interactive)
	if err != nil {
		t.Error(err)
	}
	v := make(url.Values)
	v.Set("q", "amy lane fish")
	res := a.Search(v)
	res.Print()
}

func audibleSrch(q string) []any {
	s := audible.NewSearch(q)
	r, err := audible.Products().Search(s)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("products %v\n", r.Products)
	return lo.ToAnySlice(r.Products)
}
