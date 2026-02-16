//go:build js && wasm

package main

import (
	"encoding/json"
	"lemin/core"
	"syscall/js"
)

type ResponseData struct {
	Farm  *core.Farm  `json:"data,omitempty"`
	Turns []core.Turn `json:"turns,omitempty"`
}

type Response struct {
	OK    bool         `json:"ok"`
	Data  ResponseData `json:"data,omitempty"`
	Error string       `json:"error,omitempty"`
}

func marshal(resp Response) string {
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func marshalError(e string) string {
	resp := Response{OK: false, Error: e}
	return marshal(resp)
}

func run(this js.Value, args []js.Value) any {
	if len(args) != 1 || args[0].Type() != js.TypeString {
		return marshalError("Lemin.run() requires 1 string argument")
	}

	farm, err := core.Parse(args[0].String())
	if err != nil {
		return marshalError(err.Error())
	}

	turns, err := core.Lemin(farm)
	if err != nil {
		return marshalError(err.Error())
	}

	resp := Response{OK: true, Data: ResponseData{
		Farm:  farm,
		Turns: turns,
	}}
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
