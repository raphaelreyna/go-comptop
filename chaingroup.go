package comptop

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// ChainGroup represents a group in the chain complex of a Complex.
// The elements of a chain group are the chains in a Complex c, consisting of p-dimensional simplices in c.
// In other words, a ChainGroup represents all possible combinations of the p-dimensional simplices of a Complex.
// Since ChainGroup represents a group, chains witihn a ChainGroup can be added together to obtain a new Chain in the ChainGroup.
// Every ChainGroup is a boolean group which means that the sum of Chain with itself is an empty Chain.
//
// More info: https://en.wikipedia.org/wiki/Chain_complex
type ChainGroup struct {
	complex *Complex

	simplices map[Index]*Simplex
	basespace map[Index]struct{}
	dim       Dim

	zero *Chain

	head Index

	boundaryMatrix mat.Matrix
}

func (c *Complex) newChainGroup(dim Dim) *ChainGroup {
	cg := &ChainGroup{
		complex:   c,
		simplices: map[Index]*Simplex{},
		basespace: map[Index]struct{}{},
		zero:      &Chain{complex: c, dim: dim},
		dim:       dim,
	}

	cg.zero.chaingroup = cg

	return cg
}

func (cg *ChainGroup) addSimplex(s *Simplex) {
	if s.Dim() != cg.dim {
		return
	}

	s.index = cg.head
	if _, exists := cg.simplices[s.index]; exists {
		return
	}

	cg.head++

	cg.simplices[s.index] = s
	for _, v := range s.base {
		cg.basespace[v] = struct{}{}
	}

	cg.boundaryMatrix = nil

	higherGroup := cg.complex.chainGroups[cg.dim+1]
	if higherGroup != nil {
		higherGroup.boundaryMatrix = nil
	}
}

// Dim is the dimension of the simplices that make up the ChainGroup.
func (cg *ChainGroup) Dim() Dim {
	return cg.dim
}

// Rank is the number of simplices that make up the ChainGroup.
func (cg *ChainGroup) Rank() int {
	return len(cg.simplices)
}

// Simplices returns the simplices that make up the ChainGroup.
func (cg *ChainGroup) Simplices() []*Simplex {
	els := []*Simplex{}

	if cg.simplices == nil {
		cg.simplices = map[Index]*Simplex{}
	}

	for _, s := range cg.simplices {
		els = append(els, s)
	}

	return els
}

// Simplex returns the Simplex with index idx from the ChainGroups set of simplices.
func (cg *ChainGroup) Simplex(idx Index) *Simplex {
	return cg.simplices[idx]
}

func (cg *ChainGroup) String() string {
	return fmt.Sprintf(`ChainGroup{dim: %d, rank: %d, elements: %v}`,
		cg.dim,
		cg.Rank(),
		cg.simplices,
	)
}

// IsElement returns true if c is a Chain in the ChainGroup; returns false otherwise.
func (cg *ChainGroup) IsElement(c *Chain) bool {
	if c.complex != cg.complex {
		return false
	}

	if c.dim != cg.dim {
		return false
	}

	for v := range c.base {
		if _, exists := cg.basespace[v]; !exists {
			return false
		}
	}

	return true
}

// Add returns the Chain which is a result of adding the Chain containing only a and the Chain containing only b.
func (cg *ChainGroup) Add(a, b *Simplex) *Chain {
	if a == nil && b == nil {
		return nil
	}

	if a == nil {
		return cg.Singleton(b)
	}

	if b == nil {
		return cg.Singleton(a)
	}

	if a.complex != b.complex {
		return nil
	}

	// Make sure dimensions match
	if a.Dim() != b.Dim() {
		return nil
	}
	if cg.dim != a.Dim() {
		return nil
	}

	// Return the zero element Chain if a == b
	if a.Equals(b) {
		return cg.zero
	}

	// Copy over the base elements
	base := map[Index]struct{}{}
	for _, v := range a.base {
		base[v] = struct{}{}
	}
	for _, v := range b.base {
		base[v] = struct{}{}
	}

	return &Chain{
		chain: chain{
			simplices: []*Simplex{a, b},
		},
		complex:    cg.complex,
		chaingroup: cg,
		dim:        cg.dim,
		idxs: map[Index]*Simplex{
			a.index: a,
			b.index: b,
		},
		base: base,
	}
}

// Singleton returns the Chain consisting of only s.
func (cg *ChainGroup) Singleton(s *Simplex) *Chain {
	if cg.dim != s.dim() {
		return nil
	}

	if cg.complex != s.complex {
		return nil
	}

	chain := &Chain{
		chain: chain{
			simplices: []*Simplex{s},
		},
		complex:    cg.complex,
		chaingroup: cg,
		dim:        cg.dim,
		idxs: map[Index]*Simplex{
			s.index: s,
		},
		base: map[Index]struct{}{},
	}

	for _, v := range s.base {
		chain.base[v] = struct{}{}
	}

	return chain
}

// IsZero returns true if c is the zero element of the ChainGroup.
func (cg *ChainGroup) IsZero(c *Chain) bool {
	if c.dim != cg.dim {
		return false
	}

	if c.complex != c.complex {
		return false
	}

	if len(c.simplices) == 0 {
		return true
	}

	return false
}

// Zero returns the zero element of the ChainGroup.
func (cg *ChainGroup) Zero() *Chain {
	return cg.zero
}

// ChainFromVector returns the Chain represented by v.
func (cg *ChainGroup) ChainFromVector(v Vector) *Chain {
	if r, _ := v.Dims(); r != cg.Rank() {
		return nil
	}

	chain := &Chain{
		chain: chain{
			simplices: []*Simplex{},
		},
		complex:    cg.complex,
		chaingroup: cg,
		dim:        cg.dim,
		idxs:       map[Index]*Simplex{},
		base:       map[Index]struct{}{},
	}

	vv := v.(mat.Matrix)
	vm, _ := vv.Dims()

	for idx := 0; idx < vm; idx++ {
		if vv.At(idx, 0) == 1.0 {
			i := Index(idx)
			smplx := cg.simplices[i]
			chain.simplices = append(chain.simplices, smplx)
			chain.idxs[i] = smplx
			for _, vert := range smplx.base {
				chain.base[vert] = struct{}{}
			}
		}
	}

	return chain
}

func (cg *ChainGroup) NewChainFromSimplices(s ...*Simplex) *Chain {
	chain := &Chain{
		chain: chain{
			simplices: []*Simplex{},
		},
		complex:    cg.complex,
		chaingroup: cg,
		dim:        cg.dim,
		idxs:       map[Index]*Simplex{},
		base:       map[Index]struct{}{},
	}

	for _, smplx := range s {
		if smplx.Dim() != chain.dim {
			continue
		}

		chain.simplices = append(chain.simplices, smplx)
		chain.idxs[smplx.index] = smplx
		for _, vert := range smplx.base {
			chain.base[vert] = struct{}{}
		}
	}

	return chain
}

func (cg *ChainGroup) BoundaryMatrix() mat.Matrix {
	if cg.dim == 0 {
		return nil
	}

	if cg.boundaryMatrix != nil {
		return cg.boundaryMatrix
	}

	cg.updateBoundaryMatrix()

	return cg.boundaryMatrix
}

func (cg *ChainGroup) updateBoundaryMatrix() {
	if cg.dim == 0 {
		return
	}

	complex := cg.complex
	lowerGroup := complex.chainGroups[cg.dim-1]
	if lowerGroup == nil {
		complex.chainGroups[cg.dim-1] = cg.complex.newChainGroup(cg.dim - 1)
		lowerGroup = complex.chainGroups[cg.dim-1]
	}

	n := cg.Rank()         // cols
	m := lowerGroup.Rank() // rows

	if m == 0 || n == 0 {
		return
	}

	bm := mat.NewDense(m, n, make([]float64, n*m))

	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			lowDimSimplex := lowerGroup.simplices[Index(row)]
			highDimSimplex := cg.simplices[Index(col)]

			a := 0.0
			if highDimSimplex.HasFace(lowDimSimplex) {
				a = 1.0
			}

			bm.Set(row, col, a)
		}
	}

	cg.boundaryMatrix = bm
}

type ChainGroups map[Dim]*ChainGroup
