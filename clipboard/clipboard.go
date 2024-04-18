//go:build !wasm
// +build !wasm

package clipboard

import (
	"golang.design/x/clipboard"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}

func Write(text string) {
	clipboard.Write(clipboard.FmtText, []byte(text))
}

func Read() string {
	return string(clipboard.Read(clipboard.FmtText))
}
