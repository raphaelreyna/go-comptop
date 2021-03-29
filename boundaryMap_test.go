package comptop

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestBoundaryMap_reduce(t *testing.T) {
	// The d_1 boundary map of the tetrahedron
	bm := &BoundaryMap{
		mat: mat.NewDense(4, 6, []float64{
			1, 1, 1, 0, 0, 0,
			1, 0, 0, 1, 1, 0,
			0, 1, 0, 1, 0, 1,
			0, 0, 1, 0, 1, 1,
		}),
	}

	expectedU := mat.NewDense(4, 4, []float64{
		1, 0, 0, 0,
		1, 1, 0, 0,
		1, 1, 1, 0,
		1, 1, 1, 1,
	})

	expectedV := mat.NewDense(6, 6, []float64{
		1, 1, 0, 1, 1, 0,
		0, 1, 1, 1, 0, 1,
		0, 0, 1, 0, 1, 1,
		0, 0, 0, 1, 0, 0,
		0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 1,
	})

	bm.reduce()

	if !mat.Equal(bm.u, expectedU) {
		t.Error("invalid U_0 matrix")
	}

	if !mat.Equal(bm.v, expectedV) {
		t.Error("invalid V_1 matrix")
	}

	// The d_2 boundary map of the tetrahedron
	bm = &BoundaryMap{
		mat: mat.NewDense(6, 4, []float64{
			1, 1, 0, 0,
			1, 0, 1, 0,
			0, 1, 1, 0,
			1, 0, 0, 1,
			0, 1, 0, 1,
			0, 0, 1, 1,
		}),
	}

	expectedV = mat.NewDense(4, 4, []float64{
		1, 1, 1, 1,
		0, 1, 1, 1,
		0, 0, 1, 1,
		0, 0, 0, 1,
	})

	expectedU = mat.NewDense(6, 6, []float64{
		1, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 0,
		1, 1, 1, 0, 0, 0,
		0, 1, 0, 1, 0, 0,
		1, 0, 0, 1, 1, 0,
		0, 1, 0, 1, 0, 1,
	})

	bm.reduce()

	if !mat.Equal(bm.u, expectedU) {
		t.Error("invalid U_1 matrix")
	}

	if !mat.Equal(bm.v, expectedV) {
		t.Error("invalid V_2 matrix")
	}

}
