package srch

import (
	"net/url"
	"testing"
)

func TestNewQuery(t *testing.T) {
	q := testQuery()
	println(q.String())
}

func testQuery() Query {
	vals := make(url.Values)
	vals.Add("tags", "abo")
	vals.Add("tags", "dnr")
	vals.Add("authors", "Alice Winters")
	vals.Add("authors", "Amy Lane")
	//vals.Add("q", "fish")
	return Query(vals)
}
