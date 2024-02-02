package main

import (
	"syscall/js"
)

func RegisterFunc(name string, fn func(this js.Value, args []js.Value) any) {
	js.Global().Get("srch").Set(name, js.FuncOf(fn))
}

func RegisterFuncWithNoArgs(name string, fn func() any) {
	jf := func(this js.Value, args []js.Value) any {
		return fn()
	}
	RegisterFunc(name, jf)
}

func init() {
}

func main() {

	//println("wasm loaded")
	<-make(chan bool)
}
