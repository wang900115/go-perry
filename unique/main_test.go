package main

import (
	"testing"
	"unique"
)

type structA struct {
	Content string
}

var (
	s1 = structA{"just a really long string"}
	s2 = structA{"just a really long string"}
)

func BenchmarkStringComparisong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = (s1 == s2)
	}
}

func BenchmarkUniqueHandleComparisong(b *testing.B) {

	h1 := unique.Make(s1)
	h2 := unique.Make(s2)
	for i := 0; i < b.N; i++ {
		_ = (h1 == h2)
	}
}
