//go:build !wasm
// +build !wasm

package main

import (
	"fmt"
	"os"

	"github.com/elfviewer/elfviewer/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}