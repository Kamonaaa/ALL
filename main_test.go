package main

import (
	"strings"
	"testing"
)

func BenchmarkProcess(b *testing.B) {
	input := strings.Repeat("hello world (up) , ", 1000000)

	for i := 0; i < b.N; i++ {
		tokens := tokenize(input)
		processed := process(tokens)
		_ = format(processed)
	}
}
