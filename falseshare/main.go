package main

import (
	"golang.org/x/sys/cpu"
)

// Cache line size is 64 bytes, so we can use padding to avoid false sharing between A and B.

type Counter struct {
	A int64
	B int64
}

type CustomCounter struct {
	A int64 // 64 / 8 = 8
	_ cpu.CacheLinePad
	B int64 // 64 / 8 = 8
}
