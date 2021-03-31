package comptop

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestBoundaryMap_reduce(t *testing.T) {
	// The D_0 boundary matrix of the tetrahedron
	bm := &BoundaryMap{
		mat: mat.NewDense(1, 4, []float64{
			1, 1, 1, 1,
		}),
	}

	expectedSmithNormal := mat.NewDense(1, 4, []float64{
		1, 0, 0, 0,
	})

	bm.reduce()

	if !mat.Equal(bm.sn, expectedSmithNormal) {
		fsn := mat.Formatted(bm.sn)
		t.Errorf("invalid N_0 matrix:\n%v", fsn)
	}

	a := mat.NewDense(1, 4, nil)
	a.Mul(bm.mat, bm.v)
	a.Mul(bm.u, a)

	for row := 0; row < 1; row++ {
		for col := 0; col < 4; col++ {
			a.Set(row, col, float64(int(a.At(row, col))%2))
		}
	}

	if !mat.Equal(a, bm.sn) {
		t.Error("invalid Smith-normal factorization of D_0")
	}

	// The D_1 boundary matrix of the tetrahedron
	bm = &BoundaryMap{
		mat: mat.NewDense(4, 6, []float64{
			1, 1, 1, 0, 0, 0,
			1, 0, 0, 1, 1, 0,
			0, 1, 0, 1, 0, 1,
			0, 0, 1, 0, 1, 1,
		}),
	}

	expectedSmithNormal = mat.NewDense(4, 6, []float64{
		1, 0, 0, 0, 0, 0,
		0, 1, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	})

	bm.reduce()

	if !mat.Equal(bm.sn, expectedSmithNormal) {
		fsn := mat.Formatted(bm.sn)
		t.Errorf("invalid N_1 matrix:\n%v", fsn)
	}

	a = mat.NewDense(4, 6, nil)
	a.Mul(bm.mat, bm.v)
	a.Mul(bm.u, a)

	for row := 0; row < 4; row++ {
		for col := 0; col < 6; col++ {
			a.Set(row, col, float64(int(a.At(row, col))%2))
		}
	}

	if !mat.Equal(a, bm.sn) {
		t.Error("invalid Smith-normal factorization of D_1")
	}

	// The D_2 boundary matrix of the tetrahedron
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

	expectedSmithNormal = mat.NewDense(6, 4, []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	})

	bm.reduce()

	if !mat.Equal(bm.sn, expectedSmithNormal) {
		fsn := mat.Formatted(bm.sn)
		t.Errorf("invalid N_2 matrix:\n%v", fsn)
	}

	a = mat.NewDense(6, 4, nil)
	a.Mul(bm.mat, bm.v)
	a.Mul(bm.u, a)

	for row := 0; row < 6; row++ {
		for col := 0; col < 4; col++ {
			a.Set(row, col, float64(int(a.At(row, col))%2))
		}
	}

	if !mat.Equal(a, bm.sn) {
		t.Error("invalid Smith-normal factorization of D_2")
	}

	// The D_3 boundary matrix of the tetrahedron
	bm = &BoundaryMap{
		mat: mat.NewDense(4, 1, []float64{
			1, 1, 1, 1,
		}),
	}

	expectedSmithNormal = mat.NewDense(4, 1, []float64{
		1, 0, 0, 0,
	})

	bm.reduce()

	if !mat.Equal(bm.sn, expectedSmithNormal) {
		fsn := mat.Formatted(bm.sn)
		t.Errorf("invalid N_3 matrix:\n%v", fsn)
	}

	a = mat.NewDense(4, 1, nil)
	a.Mul(bm.mat, bm.v)
	a.Mul(bm.u, a)

	for row := 0; row < 4; row++ {
		for col := 0; col < 1; col++ {
			a.Set(row, col, float64(int(a.At(row, col))%2))
		}
	}

	if !mat.Equal(a, bm.sn) {
		t.Error("invalid Smith-normal factorization of D_3")
	}
}
