package main

import (
	"fmt"
	"runtime"
	"strings"
	"unique"
)

func getUsedMem() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func main() {
	// --- A
	memStartA := getUsedMem()
	var rawStrings []string

	for i := 0; i < 10000; i++ {
		str := strings.Clone("example_string_with_some_length_12345")
		rawStrings = append(rawStrings, str)
	}
	memEndA := getUsedMem()
	runtime.KeepAlive(rawStrings)

	// --- B
	memStartB := getUsedMem()
	var handles []unique.Handle[string]

	for i := 0; i < 10000; i++ {
		str := strings.Clone("example_string_with_some_length_12345")
		handle := unique.Make(str)
		handles = append(handles, handle)
	}
	memEndB := getUsedMem()
	runtime.KeepAlive(handles)

	fmt.Printf("A: %d bytes\n", memEndA-memStartA)
	fmt.Printf("B: %d bytes\n", memEndB-memStartB)
}
