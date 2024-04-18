//go:build wasm
// +build wasm

package clipboard

import (
	"syscall/js"
)

func Write(text string) {
	setResult := make(chan struct{}, 1)
	js.Global().Get("navigator").Get("clipboard").Call("writeText", js.ValueOf(text)).
		Call("then",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				setResult <- struct{}{}
				return nil
			})).
		Call("catch",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				println("failed to set clipboard: " + args[0].String())
				setResult <- struct{}{}
				return nil
			}),
		)

	<-setResult
}

func Read() string {
	setResult := make(chan string, 1)
	js.Global().Get("navigator").Get("clipboard").Call("readText").
		Call("then",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				setResult <- args[0].String()
				return nil
			})).
		Call("catch",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				println("failed to get clipboard: " + args[0].String())
				setResult <- ""
				return nil
			}),
		)

	return <-setResult
}
