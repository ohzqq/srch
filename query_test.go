package srch

import (
	"fmt"
	"net/url"
	"testing"
)

func TestNewQuery(t *testing.T) {
	q := testQuery()
	println(q.String())

	kw := q.Keywords()
	fmt.Printf("kw %v\n", kw)

	f := q.Filters()
	fmt.Printf("filters %v\n", f)
}

func testQuery() Query {
	vals := make(url.Values)
	vals.Add("tags", "abo")
	vals.Add("tags", "dnr")
	vals.Add("authors", "Alice Winters")
	vals.Add("authors", "Amy Lane")
	vals.Add("q", "fish")
	return Query(vals)
}
