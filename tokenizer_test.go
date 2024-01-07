package srch

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

type index map[string][]int

func TestFacetSearch(t *testing.T) {
	facet := NewFacet("title")
	for _, b := range books {
		title := cast.ToString(b["title"])
		for _, token := range Tokenizer(title) {
			facet.AddItem(token, cast.ToString(b["id"]))
		}
	}
	bits := facet.Filter("fish")
	ids := bits.ToArray()
	//filtered := FilteredItems(books, lo.ToAnySlice(ids))
	fmt.Printf("%v\n", ids)
}

func TestIndexSearch(t *testing.T) {
	index := make(index)
	index.add(books)
	res := index.Search("fish")
	fmt.Printf("%v\n", res)
}

func (idx index) add(docs []map[string]any) {
	for _, b := range docs {
		title := cast.ToString(b["title"])
		id := cast.ToInt(b["id"])
		for _, token := range Tokenizer(title) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == id {
				continue
			}
			idx[token] = append(ids, id)
		}
	}
}

func (idx index) Search(text string) [][]int {
	var r [][]int
	for _, token := range Tokenizer(text) {
		if ids, ok := idx[token]; ok {
			r = append(r, ids)
		}
	}
	return r
}

var titles = []string{
	"Sporemaggeddon Vol. 1",
	"Red Fish, Dead Fish",
	"Apocalypse: Regression, Book 1",
	"The Land - Forging",
	"100th Run, Book One",
	"Fish on a Bicycle",
}
