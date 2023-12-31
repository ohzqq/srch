package srch

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
)

type testSearcher struct {
	cmd []string
}

var testS = testSearcher{
	cmd: []string{"list", "--with-library", "http://localhost:8888/#audiobooks", "--username", "churkey", "--password", "<f:/home/mxb/.dotfiles/misc/calibre.txt>", "--limit", "2", "--for-machine"},
}

type testQ string

func (q testQ) String() string {
	return string(q)
}

func TestCDB(t *testing.T) {
	s := New(testS)
	//err := s.Get()
	err := s.Get(testQ("litrpg"))
	if err != nil {
		t.Error(err)
	}
	sel, err := s.Results()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", sel)
}

func cdbSearch(t *testing.T) []byte {
	//cdb := exec.Command("echo", `angst`)

	cdb := exec.Command("calibredb", testS.cmd...)
	//println(cdb.String())

	out, err := cdb.Output()
	if err != nil {
		t.Error(err)
	}
	return out
}

type testResult map[string]any

func (s testSearcher) Search(queries ...Query) ([]Result, error) {
	if len(queries) > 0 {
		testS.cmd = append(testS.cmd, "-s", queries[0].String())
	}
	cdb := exec.Command("calibredb", testS.cmd...)
	println(cdb.String())

	out, err := cdb.Output()
	if err != nil {
		return nil, err
	}

	var res []testResult
	err = json.Unmarshal(out, &res)
	if err != nil {
		return nil, err
	}

	var items []Result
	for _, r := range res {
		items = append(items, r)
	}

	return items, nil
}

func (r testResult) String() string {
	return r["title"].(string)
}

func (s testSearcher) Find() SearchFunc {
	return s.Search
}
