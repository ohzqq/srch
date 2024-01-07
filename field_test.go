package srch

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
)

func TestFieldSearch(t *testing.T) {
	facet := NewField("title")
	facet.FieldType = Text
	for _, book := range books {
		b := book.(map[string]any)
		title := cast.ToString(b["title"])
		for _, token := range Tokenizer(title) {
			facet.Add(token, cast.ToInt(b["id"]))
		}
	}

	bits := facet.Search("fish")
	//ids := bits.ToArray()
	//filtered := FilteredItems(books, lo.ToAnySlice(ids))
	fmt.Printf("%v\n", bits)
}
