package main

import (
	"testing"
)

// BenchmarkJaggedSlice measures the performance of your specific allocation logic.
func BenchmarkAppendSlice(b *testing.B) {
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
func BenchmarkCopySlice(b *testing.B) {
	size := 100
	b.ResetTimer() // Start timing after any setup logic
	for n := 0; n < b.N; n++ {
		s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
// To delete index 5:
// We copy everything from index 6 onwards onto the position of index 5
copy(s[5:], s[6:]) 

// Resulting s: [0 1 2 3 4 6 7 8 9 9]
s = s[:len(s)-1] // Manually shrink the slice to hide the extra '9'
	}
}
