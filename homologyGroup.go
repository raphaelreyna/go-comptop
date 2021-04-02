package comptop

import (
	"gonum.org/v1/gonum/mat"
)

// CycleGroup Z_p is a subgroup of the ChainGroup C_p of the same dimension p.
// A cycle group Z_p consists of all chains in in the chain group C_p with a zero / empty boundary (ie cycles).
type CycleGroup struct {
	chainGroup *ChainGroup

	basis []*Chain
}

func (cg *ChainGroup) CycleGroup() *CycleGroup {
	if cg.cg != nil {
		return cg.cg
	}

	cg.cg = &CycleGroup{
		chainGroup: cg,
		basis:      []*Chain{},
	}

	bm := cg.BoundaryMap()
	v := bm.V()
	l := bm.SmithNormalDiagonalLength()
	rank := cg.Rank()

	for i := l; i < rank; i++ {
		col := mat.Col(nil, i, v)
		chain := cg.ChainFromVector(Vector(mat.NewDense(rank, 1, col)))
		cg.cg.basis = append(cg.cg.basis, chain)
	}

	return cg.cg
}

func (cg *CycleGroup) Basis() []*Chain {
	return cg.basis
}

func (cg *CycleGroup) ChainGroup() *ChainGroup {
	return cg.chainGroup
}

func (cg *CycleGroup) Rank() int {
	return len(cg.basis)
}

type BoundaryGroup struct {
	chainGroup *ChainGroup

	basis []*Chain
}

// BoundaryGroup B_p is a subgroup of the CycleGroup Z_p of the same dimension p.
// A boundary group B_p consists of all cycles in the cycle group Z_p which are the boundary of a chain in C_{p+1}.
func (cg *ChainGroup) BoundaryGroup() *BoundaryGroup {
	if cg.bg != nil {
		return cg.bg
	}

	cg.bg = &BoundaryGroup{
		chainGroup: cg,
		basis:      []*Chain{},
	}

	higherGroup := cg.complex.ChainGroup(cg.Dim() + 1)
	var bm *BoundaryMap
	if higherGroup != nil {
		bm = higherGroup.BoundaryMap()
	} else {
		data := make([]float64, cg.Rank())
		for i := 0; i < cg.Rank(); i++ {
			data[i] = 1.0
		}
		bm = &BoundaryMap{
			mat: mat.NewDense(cg.Rank(), 1, data),
		}
	}
	ui := bm.UInverse()
	l := bm.SmithNormalDiagonalLength()
	rank := cg.Rank()

	for i := 0; i < l; i++ {
		col := mat.Col(nil, i, ui)
		chain := cg.ChainFromVector(Vector(mat.NewDense(rank, 1, col)))
		cg.bg.basis = append(cg.bg.basis, chain)
	}

	return cg.bg
}

func (bg *BoundaryGroup) Basis() []*Chain {
	return bg.basis
}

func (bg *BoundaryGroup) ChainGroup() *ChainGroup {
	return bg.chainGroup
}

func (bg *BoundaryGroup) Rank() int {
	return len(bg.basis)
}

// HomologyGroup H_p is the quotient of Z_p and B_p: H_p = Z_p / B_p
type HomologyGroup struct {
	chainGroup *ChainGroup

	boundaryBasis []*Chain

	basis   []*Chain
	minimal bool
}

func (cg *ChainGroup) HomologyGroup() *HomologyGroup {
	if cg.hg != nil {
		return cg.hg
	}

	cg.hg = &HomologyGroup{
		chainGroup: cg,
	}

	b := cg.BoundaryGroup()
	cg.hg.boundaryBasis = b.Basis()

	return cg.hg

}

func (hg *HomologyGroup) Basis() []*Chain {
	if hg.basis != nil {
		return hg.basis
	}

	cg := hg.chainGroup

	z := cg.CycleGroup()
	zBasis := z.Basis()
	b := cg.BoundaryGroup()

	m := mat.NewDense(cg.Rank(), z.Rank(), nil)

	for idx, chain := range b.Basis() {
		col := mat.Col(nil, 0, mat.Matrix(chain.Vector()))
		m.SetCol(idx, col)
	}

	cCombos := chainCombinations(z.Rank()-b.Rank(), zBasis)

	for _, combo := range cCombos {
		for idx, chain := range combo {
			col := mat.Col(nil, 0, mat.Matrix(chain.Vector()))
			m.SetCol(b.Rank()+idx, col)
		}

		if isLI(m) {
			hg.basis = combo
			return combo
		}
	}

	return nil
}

// MinimalBasis computes a basis for the homology group that is minimal with respect to Hamming weight + length of the interesection of the chains in the basis.
// The idea is to **try** to find the smallest cycles that cycle around the holes in the most linearly independent way.
func (hg *HomologyGroup) MinimalBasis() []*Chain {
	if hg.basis != nil && hg.minimal {
		return hg.basis
	}

	cg := hg.chainGroup

	z := cg.CycleGroup()
	zBasis := z.Basis()
	b := cg.BoundaryGroup()

	m := mat.NewDense(cg.Rank(), z.Rank(), nil)

	for idx, chain := range b.Basis() {
		col := mat.Col(nil, 0, mat.Matrix(chain.Vector()))
		m.SetCol(idx, col)
	}

	cCombos := chainCombinations(z.Rank()-b.Rank(), zBasis)

	var (
		minCombo  []*Chain
		minWeight int = 1<<31 - 1
	)

	for _, combo := range cCombos {
		var (
			intersection *Chain
			w            int
		)
		for idx, chain := range combo {
			col := mat.Col(nil, 0, mat.Matrix(chain.Vector()))
			m.SetCol(b.Rank()+idx, col)
			intersection = intersection.Intersection(chain)
			w += hammingWeight(col)
		}

		ww := intersection.Len() + w

		if ww >= minWeight {
			continue
		}

		if isLI(m) {
			minWeight = ww
			minCombo = combo
		}
	}

	return minCombo
}
