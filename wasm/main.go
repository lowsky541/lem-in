//go:build js && wasm

package main

import (
	"encoding/json"
	"lemin/core"
	"syscall/js"
)

type Response struct {
	OK     bool         `json:"ok"`
	Result *core.Result `json:"data,omitempty"`
	Error  string       `json:"error,omitempty"`
}

func marshal(resp Response) string {
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func marshalError(err string, res *core.Result) string {
	resp := Response{OK: false, Error: err, Result: res}
	return marshal(resp)
}

func run(this js.Value, args []js.Value) any {
	if len(args) != 1 || args[0].Type() != js.TypeString {
		return marshalError("Lemin.run() requires 1 string argument", nil)
	}

	farmDesc := args[0].String()
	res, err := core.Run(farmDesc)
	if err != nil {
		return marshalError(err.Error(), res)
	}

	resp := Response{OK: true, Result: res}
	return marshal(resp)
}

func main() {
	js.Global().Set(
		"Lemin",
		js.ValueOf(map[string]any{
			"run": js.FuncOf(run),
		}),
	)

	// Keep the WebAssembly program running indefinitely
	// Without this, the program would exit and JavaScript couldn't call our function
	select {}
}
