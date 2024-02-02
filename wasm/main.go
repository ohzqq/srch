package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/ohzqq/srch"
)

var idx *srch.Index

func NewClient(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return nil
	}
	q := args[0].String()
	i, err := srch.New(q)
	if err != nil {
		println("client init error")
		return nil
	}
	idx = i

	return map[string]any{
		"search": js.FuncOf(Search),
	}
}

func Search(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return nil
	}
	params := args[0].String()
	requests := make(map[string][]map[string]string)
	err := json.Unmarshal([]byte(params), &requests)
	if err != nil {
		println(err.Error())
	}
	//fmt.Printf("%+v\n", requests)
	if req, ok := requests["requests"]; ok {
		if len(req) < 1 {
			return nil
		}
		p := req[0]["params"]
		res := idx.Search(p)
		results := map[string][]any{
			"results": {res.StringMap()},
		}
		d, err := json.Marshal(results)
		if err != nil {
			println(err.Error())
		}
		return string(d)
	}
	return nil
}

func init() {
	js.Global().Set("srch", make(map[string]any))
}

func main() {
	RegisterFunc("newClient", NewClient)

	println("wasm loaded")
	<-make(chan bool)
}

func RegisterFunc(name string, fn func(this js.Value, args []js.Value) any) {
	js.Global().Get("srch").Set(name, js.FuncOf(fn))
}

func RegisterFuncWithNoArgs(name string, fn func() any) {
	jf := func(this js.Value, args []js.Value) any {
		return fn()
	}
	RegisterFunc(name, jf)
}
