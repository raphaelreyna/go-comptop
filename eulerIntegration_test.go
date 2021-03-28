package comptop

import "testing"

// Simplicial Complex of sensor network
// ************************************
// Vertex Labels for simplicial complex X:
//
//  0----1----2
//   \ /   \ /
//    3-----4
//   / \   / \
//  5----6----7
//   \ /   \ /
//    8-----9
//   / \   / \
//  10---11---12
//
// Height function values h(x):
//
//  1----3----1
//   \ /   \ /
//    2-----2
//   / \   / \
//  0----0----0
//   \ /   \ /
//    0-----1
//   / \   / \
//  0----1----1
//
// Euler Integral:
//
//    /-\
//    |
//    | h(x) d\/ = 4
//    |	      /\
//  \-/
//
func TestEulerIntegral(t *testing.T) {
	c := &Complex{}
	b := []Base{
		{0, 1, 3}, {1, 3, 4}, {1, 2, 4},
		{3, 5, 6}, {3, 4, 6}, {4, 6, 7},
		{5, 6, 8}, {6, 8, 9}, {6, 7, 9},
		{8, 10, 11}, {8, 9, 11}, {9, 11, 12},
	}
	c.NewSimplices(b...)

	data := map[Index]int{
		0: 1, 1: 3, 2: 1,
		3: 2, 4: 2,
		5: 0, 6: 0, 7: 0,
		8: 0, 9: 1,
		10: 0, 11: 1, 12: 1,
	}

	f := CF(func(idx Index) int {
		return data[idx]
	})

	if a := c.EulerIntegral(0, 3, f); a != 4 {
		t.Fatalf("expected 4 targets, counted: %d", a)
	}
}
