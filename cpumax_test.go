package student

import "fmt"

// Grid represents a 2D plane stored in a 1D slice
type Grid struct {
	Rows int
	Cols int
	Data []int
}

// NewGrid allocates memory once (the fast way)
func NewGrid(rows, cols int) *Grid {
	return &Grid{
		Rows: rows,
		Cols: cols,
		Data: make([]int, rows*cols),
	}
}

// Set uses the formula: Index = (row * width) + column
func (g *Grid) Set(r, c, val int) {
	g.Data[r*g.Cols+c] = val
}

// Get retrieves the value from the flat slice
func (g *Grid) Get(r, c int) int {
	return g.Data[r*g.Cols+c]
}

func main() {
	g := NewGrid(3, 3)
	g.Set(1, 2, 42)

	fmt.Printf("Value at (1,2): %d\n", g.Get(1, 2))
}
