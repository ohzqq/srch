package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
)

func handler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	//header.Set("Content-Type", param.NdJSON+"; charset=utf-8")
	header.Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintln(w, r.URL.String())

	//r.ParseForm()
	//form := r.Form

	//dec := json.NewDecoder(r.Body)
	//defer r.Body.Close()

	//var data []map[string]any
	//err := dec.Decode(&data)
	//if err != nil {
	//log.Fatal(err)
	//}

	//idx, err := srch.New(r.URL.String())
	//if err != nil {
	//log.Fatal(err)
	//}

	//err = idx.Batch(data)
	//if err != nil {
	//log.Fatal(err)
	//}

	//res, err := idx.Search(r.URL.String())
	//if err != nil {
	//log.Fatal(err)
	//}

	//enc := json.NewEncoder(w)

	//err = enc.Encode(res)
	//if err != nil {
	//log.Fatal(err)
	//}
}

func main() {
	err := cgi.Serve(http.HandlerFunc(handler))
	if err != nil {
		fmt.Println(err)
	}
}
