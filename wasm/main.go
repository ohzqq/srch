package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"syscall/js"

	"github.com/ohzqq/srch"
)

var (
	search           js.Value
	idx              *srch.Index
	NotEnoughArgsErr = errors.New("error: not enough args")
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
	//println("new client")

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

	var params string
	params = args[0].String()
	fmt.Printf("wasm search func %v\n", params)
	println(params)

	//vals, err := url.ParseQuery(params)
	//if err != nil {
	//println(err.Error())
	//return js.Null()
	//}

	//res, err := idx.Search("?" + vals.Encode())
	res, err := idx.Search(params)
	if err != nil {
		println(err.Error())
		return js.Null()
	}

	d, err := json.Marshal(res)
	if err != nil {
		println(err.Error())
		return js.Null()
	}

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
