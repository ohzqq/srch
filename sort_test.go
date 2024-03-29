package srch

import (
	"testing"

	"github.com/ohzqq/srch/param"
	"github.com/spf13/cast"
)

func TestSortAttr(t *testing.T) {
	var sortAttrs = []string{
		"title",
		"title:string",
	}

	for _, test := range sortAttrs {
		s := NewSort(test)
		if s.Field != "title" {
			t.Errorf("got %s, expected %s\n", s.Field, "title")
		}
		if s.Type != "string" {
			t.Errorf("got %s, expected %s\n", s.Type, "string")
		}
	}
}

func TestSortByInt(t *testing.T) {
	req := NewRequest().
		SetRoute(param.Dir.String()).
		UID("id").
		Facets("tags", "authors", "narrators", "series").
		SortAttr("added_stamp:int").
		SetPath(testDataDir)

	res, err := idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}
	if title := cast.ToString(res.Hits[0]["title"]); title != "Cross & Crown" {
		t.Errorf("got %s, wanted %s", title, "Cross & Crown")
	}

	req.SortBy("added_stamp")
	res, err = idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}

	if title := cast.ToString(res.Hits[0]["title"]); title != "Camp H.O.W.L." {
		t.Errorf("got %s, wanted %s", title, "Camp H.O.W.L.")
	}

	req.SortBy("added_stamp").Order("desc")
	res, err = idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}

	if title := cast.ToString(res.Hits[0]["title"]); title != "Bonds of Blood" {
		t.Errorf("got %s, wanted %s", title, "Bonds of Blood")
	}
}

func TestSortByStr(t *testing.T) {
	req := NewRequest().
		SetRoute(param.Dir.String()).
		UID("id").
		Facets("tags", "authors", "narrators", "series").
		SetPath(testDataDir)

	res, err := idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}
	if title := cast.ToString(res.Hits[0]["title"]); title != "Cross & Crown" {
		t.Errorf("got %s, wanted %s", title, "Cross & Crown")
	}

	req.SortAttr("title:string")

	req.SortBy("title")
	res, err = idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}

	if title := cast.ToString(res.Hits[0]["title"]); title != "#Blur" {
		t.Errorf("got %s, wanted %s", title, "#Blur")
	}

	req.SortBy("title").Order("desc")
	res, err = idx.Search(req.String())
	if err != nil {
		t.Error(err)
	}

	if title := cast.ToString(res.Hits[0]["title"]); title != "‘Nother Sip of Gin" {
		t.Errorf("got %s, wanted %s", title, "‘Nother Sip of Gin")
	}
}
