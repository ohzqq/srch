package idx

import (
	"maps"
	"testing"

	"github.com/ohzqq/srch/db"
)

func TestNewData(t *testing.T) {
	data, err := NewData("")
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]string{
		"index":          "",
		"index-settings": "",
	}

	if !maps.Equal(data.Tables, want) {
		t.Errorf("got %#v tables, wanted %#v\n", data.Tables, want)
	}
}

func TestNewDataDisk(t *testing.T) {
	data, err := NewData(testHareDskDir, db.WithDisk(testHareDskDir))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"index":          testHareDskDir,
		"index-settings": testHareDskDir,
	}

	if !maps.Equal(data.Tables, want) {
		t.Errorf("got %#v tables, wanted %#v\n", data.Tables, want)
	}
}

func TestNewDataNet(t *testing.T) {
	data, err := NewData(testHareURL, db.WithDisk(testHareDskDir))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"index":          testHareDskDir,
		"index-settings": testHareDskDir,
	}

	if !maps.Equal(data.Tables, want) {
		t.Errorf("got %#v tables, wanted %#v\n", data.Tables, want)
	}
}
