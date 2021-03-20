package comptop

import "gonum.org/v1/gonum/mat"

type Vector mat.Matrix

func NewVector(els ...int) Vector {
	v := mat.NewDense(len(els), 1, nil)
	for i := range els {
		v.Set(i, 0, float64(els[i]))
	}

	return mat.Matrix(v).(Vector)
}
