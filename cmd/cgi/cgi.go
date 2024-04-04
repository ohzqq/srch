package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ohzqq/srch/param"
)

func handler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", param.NdJSON+"; charset=utf-8")

	//r.ParseForm()
	//form := r.Form

	dec, err := json.NewDecoder(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	var data []map[string]any
	err = dec.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	idx, err := New(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}

	err = idx.Batch(data)
	if err != nil {
		log.Fatal(err)
	}

	res, err := idx.Search(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}

	enc, err := json.NewEncoder(w)
	if err != nil {
		log.Fatal(err)
	}

	err = enc.Encode(res)
	if err != nil {
		log.Fatal(err)
	}
}
