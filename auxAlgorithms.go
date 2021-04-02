package comptop

import (
	"gonum.org/v1/gonum/mat"
)

// chainCombinations computes all subsets with k elements of l
func chainCombinations(k int, l []*Chain) [][]*Chain {
	combos := [][]*Chain{}

	if k >= len(l) || k <= 0 {
		return combos
	}

	n := len(l)

	for start := 0; start < n-k+1; start++ {
		end := start + k - 1

		head := make([]*Chain, k-1)
		copy(head, l[start:end])

		for idx := end; idx < n; idx++ {
			tail := l[idx]
			combos = append(combos, append(head, tail))
		}
	}

	return combos
}

// rref computes the reduced row echelon form of m.
// This is mostly stolen from https://rosettacode.org/wiki/Reduced_row_echelon_form#Go
func rref(m *mat.Dense) {
	_, cols := m.Dims()
	ele := m.RawMatrix().Data

	lead := 0
	for rxc0 := 0; rxc0 < len(ele); rxc0 += cols {
		if lead >= cols {
			return
		}
		ixc0 := rxc0
		for ele[ixc0+lead] == 0 {
			ixc0 += cols
			if ixc0 == len(ele) {
				ixc0 = rxc0
				lead++
				if lead == cols {
					return
				}
			}
		}
		for c, ix, rx := 0, ixc0, rxc0; c < cols; c++ {
			ele[ix], ele[rx] = ele[rx], ele[ix]
			ix++
			rx++
		}
		if d := ele[rxc0+lead]; d != 0 {
			d := 1 / d
			for c, rx := 0, rxc0; c < cols; c++ {
				ele[rx] *= d
				rx++
			}
		}
		for ixc0 = 0; ixc0 < len(ele); ixc0 += cols {
			if ixc0 != rxc0 {
				f := ele[ixc0+lead]
				for c, ix, rx := 0, ixc0, rxc0; c < cols; c++ {
					ele[ix] -= ele[rx] * f
					ix++
					rx++
				}
			}
		}
		lead++
	}
}

// isLi returns true is the columns of m are linearly independent, false otherwise.
func isLI(m *mat.Dense) bool {
	rref(m)

	rows, cols := m.Dims()

	for col := 0; col < cols; col++ {
	rowLoop:
		for row := 0; row < rows; row++ {
			switch m.At(row, col) {
			case 0:
				continue
			case 1:
				break rowLoop
			default:
				return false
			}
		}
	}

	return true
}

func hammingWeight(v []float64) int {
	var w int

	for _, vv := range v {
		if vv != 0 {
			w++
		}
	}

	return w
}
