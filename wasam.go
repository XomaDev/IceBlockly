//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"strings"
	"syscall/js"
)

func safeExec(fn func() js.Value) js.Value {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()
	return fn()
}

func xmlToMist(this js.Value, p []js.Value) any {
	return safeExec(func() js.Value {
		if len(p) < 1 {
			return js.ValueOf("No XML content provided")
		}
		xmlContent := p[0].String()
		groups := ParseBlockly(xmlContent)
		var builder strings.Builder

		for _, group := range groups {
			for _, block := range group {
				builder.WriteString(block.String())
				if block.Order() > 0 {
					builder.WriteString("\n")
				}
			}
			builder.WriteString("\n")
		}
		return js.ValueOf(builder.String())
	})
}

func main() {
	fmt.Println("Hello from wasam.go!")

	c := make(chan struct{}, 0)
	js.Global().Set("xmlToMist", js.FuncOf(xmlToMist))
	<-c
}
