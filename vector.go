package comptop

import "gonum.org/v1/gonum/mat"

// Vector is a vector representation of the elements of a ChainGroup of rank p where p = length of the vector.
// All elements/entries are expected to be 0 or 1.
// A 1 in the i_th position indicates that the p-dimensional Simplex with index i is part of the chain.
type Vector mat.Matrix

// NewVector returns a vector with elements/entries els.
func NewVector(els ...int) Vector {
	v := mat.NewDense(len(els), 1, nil)
	for i := range els {
		v.Set(i, 0, float64(els[i]))
	}

	return mat.Matrix(v).(Vector)
}
