//go:build wasm
// +build wasm

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"github.com/elfviewer/elfviewer/elf"
)

type ELFInfo struct {
	Ident          elf.Ident           `json:"ident"`
	Class          uint8               `json:"class"`
	Type           uint16              `json:"type"`
	Machine        uint16              `json:"machine"`
	Entry          uint64              `json:"entry"`
	SectionHeaders []elf.SectionHeader `json:"sectionHeaders"`
	ProgramHeaders []elf.ProgramHeader `json:"programHeaders"`
	Symbols        []elf.Symbol        `json:"symbols"`
}

func parseELF(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return map[string]interface{}{
			"error": "Expected 1 argument",
		}
	}

	// Convert JS ArrayBuffer to Go byte slice
	arrayBuffer := args[0]
	uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)
	length := uint8Array.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, uint8Array)

	// Parse ELF
	elfFile, err := elf.Parse(data)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	// Convert to JSON-serializable structure
	info := ELFInfo{
		Ident:          elfFile.Ident,
		Class:          elfFile.Class,
		Type:           elfFile.Type,
		Machine:        elfFile.Machine,
		Entry:          elfFile.Entry,
		SectionHeaders: elfFile.SectionHeaders,
		ProgramHeaders: elfFile.ProgramHeaders,
		Symbols:        elfFile.Symbols,
	}

	result, err := json.Marshal(info)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"data": string(result),
	}
}

func getHexDump(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return map[string]interface{}{
			"error": "Expected 2 arguments: data, sectionName",
		}
	}

	// Convert ArrayBuffer to byte slice
	arrayBuffer := args[0]
	uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)
	length := uint8Array.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, uint8Array)

	sectionName := args[1].String()

	// Parse ELF
	elfFile, err := elf.Parse(data)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	// Get hex dump using a string writer
	var buf strings.Builder
	if err := elfFile.DisplayHexDump(&buf, sectionName); err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"data": buf.String(),
	}
}

func main() {
	js.Global().Set("parseELF", js.FuncOf(parseELF))
	js.Global().Set("getHexDump", js.FuncOf(getHexDump))

	// Keep the Go program running
	select {}
}