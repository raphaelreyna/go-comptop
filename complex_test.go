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
