package comptop

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

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

func TestComplex_HomologyGroup(t *testing.T) {
	c := &Complex{}
	c.NewSimplices([]Base{
		{0, 1, 3},
		{1, 3, 4},
		{1, 2, 4},
		{2, 4, 5},
		{0, 2, 5},
		{0, 3, 5},

		{3, 4, 6},
		{4, 6, 7},
		{4, 5, 7},
		{5, 7, 8},
		{3, 5, 8},
		{3, 6, 8},

		{0, 6, 7},
		{0, 1, 7},
		{1, 7, 8},
		{1, 2, 8},
		{2, 6, 8},
		{0, 2, 6},
	}...)

	g1 := c.ChainGroup(1)
	bm1 := g1.BoundaryMap()
	bm1.reduce()
	bm1.u.Inverse(bm1.u)

	g2 := c.ChainGroup(2)
	bm2 := g2.BoundaryMap()
	bm2.reduce()
	bm2.u.Inverse(bm2.u)

	z1Basis := [][]float64{}
	for i := bm1.SmithNormalDiagonalLength(); i < g1.Rank(); i++ {
		z1Basis = append(z1Basis, mat.Col(nil, i, bm1.v))
	}

	b1Basis := [][]float64{}
	for i := 0; i < bm1.SmithNormalDiagonalLength(); i++ {
		col := mat.Col(nil, i, bm2.u)
		for idx := range col {
			col[idx] = float64(int(col[idx]) % 2)
			col[idx] *= col[idx]
		}
		b1Basis = append(b1Basis, col)
	}

	hgBasis := c.ChainGroup(1).HomologyGroup().MinimalBasis()
	t.Fatal(hgBasis)

	if len(hgBasis) != 2 {
		t.Fatalf("Expected H_1(T^2) to have 2 generators, computed %d\n", len(hgBasis))
	}
}
