package comptop

import "testing"

func TestEulerChar_Chain(t *testing.T) {
	c := &Complex{}
	b := []Base{
		{0, 1, 4},
		{0, 1, 2},
		{0, 2, 3},
		{4, 6},
		{7},
	}
	c.NewSimplices(b...)

	edges := c.NewSimplices(
		Base{0, 1}, Base{1, 4},
		Base{2, 3},
	)

	cg1 := c.GetChainGroup(1)
	chain := cg1.NewChainFromSimplices(edges.Slice()...)

	if x := chain.EulerChar(); x != 2 {
		t.Fatalf("expected Euler char of %v to be 2, got %d", chain, x)
	}
}

func TestEulerChar_Complex(t *testing.T) {
	torus := []Base{
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
	}

	c := &Complex{}

	c.NewSimplices(torus...)

	if x := c.EulerChar(); x != 0 {
		t.Fatalf("expected Euler char of torus to be 0, got %d", x)
	}
}

func TestEulerChar_SimplicialSet(t *testing.T) {
	c := &Complex{}
	b := []Base{
		{0, 1, 4},
		{0, 1, 2},
		{0, 2, 3},
		{4, 6},
		{7, 8, 9},
		{6, 9},
	}
	c.NewSimplices(b...)

	ss := NewSimplicialSet(
		c.GetSimplex(0, 1, 4),
		c.GetSimplex(7, 8, 9),
		c.GetSimplex(6, 9),
	)

	if x := ss.EulerChar(); x != 2 {
		t.Fatalf("expected Euler char of %v to be 2, got %d", ss, x)
	}
}
