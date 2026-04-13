package main

import (
	"testing"
)

const Rows = 100
const Cols = 100

// BenchmarkJagged creates a slice of slices (multiple allocations)
func BenchmarkJagged(b *testing.B) {
	for n := 0; n < b.N; n++ {
		twoD := make([][]int, Rows)
		for i := 0; i < Rows; i++ {
			twoD[i] = make([]int, Cols)
			for j := 0; j < Cols; j++ {
				twoD[i][j] = i + j
			}
		}
	}
}

// BenchmarkFlat creates one single long slice (one allocation)
func BenchmarkFlat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		flat := make([]int, Rows*Cols)
		for i := 0; i < Rows; i++ {
			for j := 0; j < Cols; j++ {
				// Manual 2D mapping: row * width + col
				flat[i*Cols+j] = i + j
			}
		}
	}
}