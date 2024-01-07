package srch

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

func TestFieldSearch(t *testing.T) {
	facet := NewTextField("title")
	for _, b := range books {
		title := cast.ToString(b["title"])
		for _, token := range Tokenizer(title) {
			facet.Add(token, b["id"])
		}
	}

	bits := facet.Search("fish")
	//ids := bits.ToArray()
	//filtered := FilteredItems(books, lo.ToAnySlice(ids))
	fmt.Printf("%v\n", bits)
}
