package comptop

import "gonum.org/v1/gonum/mat"

// CF represents the group of constructible functions on the vertex set of a complex.
//
// More info: https://en.wikipedia.org/wiki/Constructible_function
type CF func(Index) int

type Dim uint

// Index is used to establiish a total order amongst simplices of the same dimension.
// Index is also used to uniquely identify simplicies up to dimension.
type Index uint

// Base is a collection of indices for 0-dimensional simplices.
type Base []Index

func (b Base) Len() int {
	return len(b)
}

func (b Base) Less(i, j int) bool {
	return b[i] < b[j]
}

func (b Base) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

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
