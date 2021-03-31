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
	c := &Complex{}
	c.NewSimplices([]Base{
		{0, 1}, {0, 2}, {1, 3}, {2, 3}, {0, 1, 2},
	}...)

	expectedBN := []int{0, 1, 0}

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
