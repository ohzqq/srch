package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"syscall/js"

	"github.com/ohzqq/srch"
)

var (
	search           js.Value
	idx              *srch.Index
	NotEnoughArgsErr = errors.New("error: not enough args")

	JSON = &jsonGlobal{
		fn: js.Global().Get("JSON"),
	}
)

func init() {
	defer Recover()
	js.Global().Set("srch", map[string]any{})
	search = js.Global().Get("srch")
	println("search initialized")
}

func main() {
	RegisterFunc("newClient", NewClient)
	RegisterFunc("search", Search)

	<-make(chan bool)
}

func NewClient(this js.Value, args []js.Value) any {
	defer Recover()
	println("new client")

	err := CheckArgs(args)
	if err != nil {
		println(err.Error())
		return js.Null()
	}

	params := args[0].String()
	//println(params)

	dstr := args[1].String()
	var data []map[string]any
	err = json.Unmarshal([]byte(dstr), &data)
	if err != nil {
		println(err.Error())
		return js.Null()
	}

	idx, err = srch.Mem(params, data)
	if err != nil {
		println(err.Error())
		return js.Null()
	}

	return nil
}

func Search(this js.Value, args []js.Value) any {
	defer Recover()

	err := CheckArgs(args)
	if err != nil {
		println(err.Error())
		return js.Null()
	}

	//v := args[0].Index(0).Get("params")
	//q := js.Global().Get("URLSearchParams").New(v).Call("toString")

	var params string
	params = args[0].String()
	//params = JSON.Stringify(args[0])
	println(params)

	vals, err := url.ParseQuery(params)
	if err != nil {
		println(err.Error())
		return js.Null()
	}
	vals.Set("facets", "tags")
	vals.Set("hitsPerPage", "25")
	//vals.Set("query", "fish")
	//println(vals.Encode())

	//tq := `?facets=authors&facets=tags&facets=narrators&facets=series&hitsPerPage=25&order=desc&searchableAttributes=title&sortBy=added&uid=id&query=fish`

	res, err := idx.Search("?" + vals.Encode())
	if err != nil {
		println(err.Error())
		return js.Null()
	}
	//res.HitsPerPage = 25

	d, err := json.Marshal(res)
	if err != nil {
		println(err.Error())
		return js.Null()
	}

	//println(string(d))
	fmt.Printf("total hits %d\n", res.NbHits)

	return string(d)
}

func RegisterFunc(name string, fn func(this js.Value, args []js.Value) any) {
	registerFunc(name, js.FuncOf(fn))
}

func registerFunc(name string, fn js.Func) {
	defer Recover()
	search.Set(name, fn)
}

func Recover() {
	if r := recover(); r != nil {
		pc, f, l, _ := runtime.Caller(3)
		msg := fmt.Sprintf(errMsg, f, l, runtime.FuncForPC(pc).Name(), r)
		println(msg)
	}
}

var errMsg = `
file: %v
  line: %d
	func: %s
	err: %w
`

func CheckArgs(args []js.Value) error {
	if len(args) < 1 {
		return NotEnoughArgsErr
	}
	return nil
}

type jsonGlobal struct {
	fn js.Value
}

func (v *jsonGlobal) Stringify(obj js.Value) string {
	defer Recover()

	if obj.Truthy() {
		return v.fn.Call("stringify", obj).String()
	}

	return ""
}

func (v *jsonGlobal) Parse(val string) js.Value {
	defer Recover()

	if val == "" {
		val = "{}"
	}

	return v.fn.Call("parse", val)
}
