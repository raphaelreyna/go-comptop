package comptop

import "testing"

func TestComplex_PrincipleSimplices(t *testing.T) {
	cmplx := &Complex{}
	cps := cmplx.NewSimplices([]Base{
		{0, 1}, {1, 2, 3, 4}, {2, 3, 4, 5}, {10},
	}...)

	cmplx.NewSimplex(1, 2, 3)

	tps := cmplx.PrincipleSimplices()

	if count := tps.Card(); count != 4 {
		t.Fatalf("expected 4 principle simplices, got %d", count)
	}

	for _, smplx := range tps.Slice() {
		if _, exists := cps.set[smplx]; !exists {
			t.Fatalf("received non-principle simplex as principle: %v", smplx)
		}
	}
}

func TestComplex_BettiNumbers(t *testing.T) {
	// The torus = S^1 x S^1 has betti numbers 1, 2, 1.
	c := &Complex{}
	c.NewSimplices([]Base{
		{0, 1, 4},
		{1, 4, 5},
		{1, 2, 5},
		{2, 5, 6},
		{0, 2, 6},
		{0, 4, 6},
		{4, 5, 7},
		{5, 7, 8},
		{5, 6, 8},
		{6, 8, 9},
		{4, 6, 9},
		{4, 7, 9},
		{0, 7, 8},
		{0, 1, 8},
		{1, 8, 9},
		{1, 2, 9},
		{2, 7, 9},
		{0, 2, 7},
	}...)

	expectedBN := []int{1, 2, 1}

	bn := c.BettiNumbers()

	if len(expectedBN) != len(bn) {
		t.Fatalf("invalid number of Betti numbers")
	}

	for idx, ebn := range expectedBN {
		if bn[idx] != ebn {
			t.Fatalf("Betti number %d is wrong; expected %d received %d", idx, ebn, bn[idx])
		}
	}
}
