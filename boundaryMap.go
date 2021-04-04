package comptop

import (
	"gonum.org/v1/gonum/mat"
)

// BoundaryMap is the bonudary map between chain groups of dimensions p and p-1.
type BoundaryMap struct {
	mat mat.Matrix

	sn *mat.Dense
	u  *mat.Dense
	ui *mat.Dense
	v  *mat.Dense

	dl  *int
	zp  *int
	bpl *int
}

// BoundaryMatrix returns the matrix representation of the boundary map.
func (bm *BoundaryMap) BoundaryMatrix() mat.Matrix {
	return bm.mat
}

// SmithNormal returns the Smith normal form of the boundary matrix.
func (bm *BoundaryMap) SmithNormal() mat.Matrix {
	if bm.sn != nil {
		return bm.sn
	}

	bm.reduce()

	return bm.sn
}

// U returns the the left-side matrix in the Smith normal factorization of the boundaty matrix.
func (bm *BoundaryMap) U() mat.Matrix {
	if bm.u == nil {
		bm.reduce()
	}

	return bm.u
}

// UInverse returns the inverse of the left-side matrix in the Smith normal factorization of the boundaty matrix.
func (bm *BoundaryMap) UInverse() mat.Matrix {
	if bm.ui != nil {
		return bm.ui
	}

	u := bm.U()
	m, n := u.Dims()
	bm.ui = mat.NewDense(m, n, nil)
	bm.ui.Inverse(u)

	// Map into into M_{n,m}(Z_2)
	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			x := int(bm.ui.At(row, col))
			x = x % 2
			x *= x

			bm.ui.Set(row, col, float64(x))
		}
	}

	return bm.ui
}

// V returns the the right-side matrix in the Smith normal factorization of the boundaty matrix.
func (bm *BoundaryMap) V() mat.Matrix {
	if bm.v == nil {
		bm.reduce()
	}

	return bm.v
}

// SmithNormalDiagonalLength returns the number of non-zero elements on the diagonal of the Smith normal form of the boundary matrix.
func (bm *BoundaryMap) SmithNormalDiagonalLength() int {
	if bm.dl != nil {
		return *bm.dl
	}

	if bm.sn == nil {
		bm.reduce()
	}

	_, n := bm.v.Dims()
	col := 0
	for ; col < n; col++ {
		empty := true
		for _, el := range mat.Col(nil, col, bm.sn) {
			if el != 0 {
				empty = false
				break
			}
		}

		if empty {
			break
		}
	}

	bm.dl = &col

	return col
}

// Zp returns the number of 0-columns in the Smith normal form of the boundary matrix.
// The value returned by Zp coincides with the rank of the kernel of the boundary matrix,
// which is the same as the rank of the cycle group Z_p < C_p.
func (bm *BoundaryMap) Zp() int {
	if bm.zp != nil {
		return *bm.zp
	}

	_, n := bm.SmithNormal().Dims()

	z := n - bm.SmithNormalDiagonalLength()
	bm.zp = &z

	return z
}

// BpLow returns the number of non-zero rows in the Smith normal form of the boundary matrix.
// The value returned by BpLow coincides with the rank of the image of the boundary matrix,
// which is the same as the rank of the boundary group B_{p-1} < Z_{p-1} < C_{p-1}.
func (bm *BoundaryMap) BpLow() int {
	if bm.bpl != nil {
		return *bm.zp
	}

	bpl := bm.SmithNormalDiagonalLength()
	bm.bpl = &bpl

	return *bm.bpl
}

// reduce computes the Smith normal form (u and v matrices are also computed along the way).
func (bm *BoundaryMap) reduce() {
	if bm.sn != nil || bm.mat == nil {
		return
	}

	// Create a copy of the boundary matrix that we can mutate into its Smith normal form
	bm.sn = mat.DenseCopyOf(bm.mat)

	// Create the U and V factors: SmithNormal = U * BoundaryMatrix * V
	m, n := bm.mat.Dims()
	bm.u = mat.NewDense(m, m, nil)
	for i := 0; i < m; i++ {
		bm.u.Set(i, i, 1.0)
	}

	bm.v = mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		bm.v.Set(i, i, 1.0)
	}

	// Carry out the Smith normal reduction
	_reduce(0, bm.sn, bm.u, bm.v)

	// Map matrices U and V into M(Z_2)
	for row := 0; row < m; row++ {
		for col := 0; col < m; col++ {
			bm.u.Set(row, col,
				float64(int(bm.u.At(row, col))%2),
			)
		}
	}

	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			bm.v.Set(row, col,
				float64(int(bm.v.At(row, col))%2),
			)
		}
	}
}

// _reduce is the actual reduction algorithm for computhing the Smith normal form of the boundary matrix.
// Taken from 'Computational Topology: An Introduction' by Edelsbrunner & Harer, pg 88.
func _reduce(x int, n, u, v *mat.Dense) {
	uDim, vDim := n.Dims()

	var doWork bool

outer:
	for rIdx := x; rIdx < uDim; rIdx++ {
		for cIdx := x; cIdx < vDim; cIdx++ {
			if n.At(rIdx, cIdx) == 1 {
				uu := mat.NewDense(uDim, uDim, nil)
				for i := 0; i < uDim; i++ {
					uu.Set(i, i, 1.0)
				}
				_swapRows(rIdx, x, n, uu)
				u.Mul(uu, u)
				vv := mat.NewDense(vDim, vDim, nil)
				for i := 0; i < vDim; i++ {
					vv.Set(i, i, 1.0)
				}
				_swapCols(cIdx, x, n, vv)
				v.Mul(v, vv)

				doWork = true

				break outer
			}
		}
	}

	if !doWork {
		return
	}

	for i := x + 1; i < uDim; i++ {
		if n.At(i, x) == 1 {
			uu := mat.NewDense(uDim, uDim, nil)
			for i := 0; i < uDim; i++ {
				uu.Set(i, i, 1.0)
			}
			_addRows(x, i, n, uu)
			u.Mul(uu, u)
		}
	}

	for i := x + 1; i < vDim; i++ {
		if n.At(x, i) == 1 {
			vv := mat.NewDense(vDim, vDim, nil)
			for i := 0; i < vDim; i++ {
				vv.Set(i, i, 1.0)
			}
			_addCols(x, i, n, vv)
			v.Mul(v, vv)
		}
	}

	_reduce(x+1, n, u, v)
}

// Below are the matrix operations layed out in
// 'Computational Topology: An Introduction' by Edelsbrunner & Harer, pgs 86-87.

func _swapRows(i, j int, a, u *mat.Dense) {
	if i == j {
		return
	}

	iRow := mat.Row(nil, i, a)
	a.SetRow(i, mat.Row(nil, j, a))
	a.SetRow(j, iRow)

	u.Set(i, j, 1.0)
	u.Set(j, i, 1.0)
	u.Set(i, i, 0.0)
	u.Set(j, j, 0.0)
}

func _swapCols(i, j int, a, v *mat.Dense) {
	if i == j {
		return
	}

	iCol := mat.Col(nil, i, a)
	a.SetCol(i, mat.Col(nil, j, a))
	a.SetCol(j, iCol)

	v.Set(i, j, 1.0)
	v.Set(j, i, 1.0)
	v.Set(i, i, 0.0)
	v.Set(j, j, 0.0)
}

// _addRows adds row l to row k ( k -> k+l )
func _addRows(l, k int, a, u *mat.Dense) {
	_, n := a.Dims()

	for col := 0; col < n; col++ {
		a.Set(k, col, float64(int(a.At(k, col)+a.At(l, col))%2))
	}

	u.Set(k, l, 1.0)
	u.Set(k, k, 1.0)
	u.Set(l, l, 1.0)
}

// _addCols adds col k to col l ( l -> k+l )
func _addCols(k, l int, a, v *mat.Dense) {
	m, _ := a.Dims()

	for row := 0; row < m; row++ {
		a.Set(row, l, float64(int(a.At(row, l)+a.At(row, k))%2))
	}

	v.Set(k, l, 1.0)
	v.Set(k, k, 1.0)
	v.Set(l, l, 1.0)
}
