package main

import (
	"testing"
)

// BenchmarkJaggedSlice measures the performance of your specific allocation logic.
func BenchmarkJaggedSlice(b *testing.B) {
	// b.N is a value managed by the Go testing tool that increases
	// until the benchmark result is stable.
	for n := 0; n < b.N; n++ {
		// The logic you want to measure
		twoD := make([][]int, 3)
		for i := 0; i < 3; i++ {
			innerLen := i + 1
			twoD[i] = make([]int, innerLen)
			for j := 0; j < innerLen; j++ {
				twoD[i][j] = i + j
			}
		}
	}
}

// BenchmarkLargeJaggedSlice tests how the logic scales with more rows.
func BenchmarkLargeJaggedSlice(b *testing.B) {
	size := 100
	b.ResetTimer() // Start timing after any setup logic
	for n := 0; n < b.N; n++ {
		twoD := make([][]int, size)
		for i := 0; i < size; i++ {
			twoD[i] = make([]int, i+1)
			for j := 0; j < i+1; j++ {
				twoD[i][j] = i + j
			}
		}
	}
}