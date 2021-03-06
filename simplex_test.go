package comptop

import (
	"fmt"
	"testing"
)

func TestSimplex_Equals(t *testing.T) {
	a := &simplex{
		base:   []Index{0, 1, 2},
		sorted: true,
	}

	// Test Reflexivity
	if !a.equals(a) {
		t.Error("failed reflexivity")
	}

	// Test Symmetry
	b := &simplex{
		base:   []Index{3, 4, 5},
		sorted: true,
	}
	if a.equals(b) != b.equals(a) {
		t.Error("failed negative symmetry")
	}
	b = &simplex{
		base:   []Index{0, 1, 2},
		sorted: true,
	}
	if a.equals(b) != b.equals(a) {
		t.Error("failed positive symmetry")
	}

	// Test positive case
	if !a.equals(b) {
		t.Error("failed positive case")
	}

	// Test negative case
	b = &simplex{
		base:   []Index{11, 2, 3},
		sorted: true,
	}
	if a.equals(b) {
		t.Error("failed negative case")
	}
}

func TestSimplex_D(t *testing.T) {
	type testcase struct {
		simplex  *simplex
		boundary []*simplex
	}

	tests := []testcase{
		{
			&simplex{[]Index{0, 1}, true},
			[]*simplex{
				{[]Index{0}, true},
				{[]Index{1}, true},
			},
		},
		{
			&simplex{[]Index{0, 1, 2}, true},
			[]*simplex{
				{[]Index{0, 1}, true},
				{[]Index{1, 2}, true},
				{[]Index{0, 2}, true},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d-Simplex", test.simplex.dim()), func(tt *testing.T) {
			computedBoundary := test.simplex.d()
			if len(computedBoundary) != len(test.boundary) {
				t.Errorf("expected %d simplices in boundary: received %d",
					len(test.boundary), len(computedBoundary),
				)
			}

			for _, s := range computedBoundary {
				found := false
				for _, ss := range test.boundary {
					if ss.equals(s) {
						found = true
					}
				}

				if !found {
					t.Errorf("incorrect boundary:\n\texpected %+v\n\treceived: %+v", test.boundary, computedBoundary)
				}
			}
		})
	}
}

func TestNewSimplex(t *testing.T) {
	cmplx := &Complex{}

	smplx := cmplx.NewSimplex(0, 1, 2)

	if smplx == nil {
		t.Fatal("received nil Simplex")
	}

	if smplx.index != 0 {
		t.Fatalf("expected 0 index for first simplex, received: %d", smplx.index)
	}

	if !smplx.simplex.equals(&simplex{base: []Index{0, 1, 2}}) {
		t.FailNow()
	}
}

func TestSimplex_Intersection(t *testing.T) {
	cmplx := &Complex{}

	smplx := cmplx.NewSimplex(0, 1, 2)
	smplx2 := cmplx.NewSimplex(1, 2, 3)

	intersection := smplx.Intersection(smplx2)
	if intersection == nil {
		t.FailNow()
	}
	if !intersection.simplex.equals(&simplex{base: []Index{1, 2}}) {
		t.FailNow()
	}
}

func TestSimplex_HasFace(t *testing.T) {
	cmplx := &Complex{}

	smplx := cmplx.NewSimplex(0, 1, 2)
	face := cmplx.GetSimplex(1, 2)

	if !smplx.HasFace(face) {
		t.Fatalf("failed positive case")
	}

	face = cmplx.NewSimplex(3, 6)
	if smplx.HasFace(face) {
		t.Fatalf("failed negative case")
	}
}

func TestSimplex_Cofaces(t *testing.T) {
	cmplx := &Complex{}

	cmplx.NewSimplices([]Base{
		{0, 1, 2}, {0, 1, 3}, {0, 1, 5}, {1, 2, 5},
		{0, 1, 10, 11},
	}...)

	smplx := cmplx.GetSimplex(0, 1)

	if smplx.Cofaces(1) != nil {
		t.Error("expected no cofaces")
	}

	cf := smplx.Cofaces(2)
	if count := cf.Card(); count != 5 {
		t.Errorf("expected 3 cofaces, got %d", count)
	}

	cf = smplx.Cofaces(3)
	if count := cf.Card(); count != 1 {
		t.Errorf("expected 1 coface, got %d", count)
	}

	cf = smplx.AllCofaces()
	if count := cf.Card(); count != 6 {
		t.Errorf("expected 6 cofaces, got %d", count)
	}
}

func TestSimplex_Boundary(t *testing.T) {
	cmplx := &Complex{}

	smplx := cmplx.NewSimplex(1, 2, 3)

	bndry := smplx.Boundary()

	if count := len(bndry.Simplices()); count != 3 {
		t.Errorf("expected 3 faces in boundary, got %d", count)
	}

	if bb := bndry.Boundary(); !bb.IsZero() {
		t.Error("failed the fundamental lemma of homology")
	}
}
