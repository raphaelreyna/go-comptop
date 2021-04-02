package comptop

import (
	"sort"

	"gonum.org/v1/gonum/mat"
)

type chain struct {
	simplices []*Simplex
	sorted    bool
}

// Len is used to satisfy the sort.Interface interface
func (c *Chain) Len() int {
	return len(c.simplices)
}

// Less is used to satisfy the sort.Interface interface
func (c *Chain) Less(i, j int) bool {
	return c.simplices[i].index < c.simplices[j].index
}

// Swap is used to satisfy the sort.Interface interface
func (c *Chain) Swap(i, j int) {
	c.simplices[i], c.simplices[j] = c.simplices[j], c.simplices[i]
}

// Chain is an element of a ChainGroup.
// Chains are formal sums over the p-dimensional Simplices of a Complex with coefficients in Z_2 = Z/2Z.
// This means that adding a Chain to itself results in an empty Chain (the zero element of the ChainGroup).
//
// More info: https://en.wikipedia.org/wiki/Simplicial_homology#Chains
type Chain struct {
	chain
	complex    *Complex
	chaingroup *ChainGroup

	vector Vector
	base   map[Index]struct{}
	idxs   map[Index]*Simplex
	dim    Dim

	eulerChar *int

	isCycle bool
}

func (c *Chain) sort() {
	sort.Sort(c)
	c.sorted = true
}

func (c *Chain) String() string {
	s := "Chain{"

	for _, smplx := range c.simplices {
		s += smplx.String() + ", "

	}

	s += "}"

	return s
}

// Dim is the dimension od the simplices in the chain.
// (all simplices in a chain have the same dimension)
func (c *Chain) Dim() Dim {
	return c.dim
}

func (c *Chain) ChainGroup() *ChainGroup {
	return c.ChainGroup()
}

func (c *Chain) IsZero() bool {
	return len(c.simplices) == 0
}

// Add returns the results of adding Chain c to Chain a.
// Since Chain is an element of a boolean group, if c == a then the resulting Chain is empty.
func (c *Chain) Add(a *Chain) *Chain {
	if a == nil {
		return c
	}

	if a.dim != c.dim {
		return nil
	}

	if c.chaingroup != a.chaingroup {
		return nil
	}

	// Count all of the simplices in the chain both chains
	simplexCount := map[*Simplex]uint{}
	for _, smplx := range a.simplices {
		simplexCount[smplx]++
	}
	for _, smplx := range c.simplices {
		simplexCount[smplx]++
	}

	chain := &Chain{
		chain: chain{
			simplices: []*Simplex{},
		},
		complex:    c.complex,
		chaingroup: c.chaingroup,
		dim:        c.dim,
		idxs:       map[Index]*Simplex{},
		base:       map[Index]struct{}{},
	}

	// only keep the simplices that show up an odd number of times
	for smplx, count := range simplexCount {
		if count%2 == 0 {
			continue
		}

		chain.simplices = append(chain.simplices, smplx)
		chain.idxs[smplx.index] = smplx
		for _, v := range smplx.base {
			chain.base[v] = struct{}{}
		}
	}

	if len(chain.simplices) == 0 {
		return c.chaingroup.zero
	}

	return chain
}

func (c *Chain) Intersection(a *Chain) *Chain {
	if a == nil {
		return c
	}

	if c == nil {
		return a
	}

	if a.dim != c.dim {
		return nil
	}

	if c.chaingroup != a.chaingroup {
		return nil
	}

	// Count all of the simplices in the chain both chains
	simplexCount := map[*Simplex]uint{}
	for _, smplx := range a.simplices {
		simplexCount[smplx]++
	}
	for _, smplx := range c.simplices {
		simplexCount[smplx]++
	}

	chain := &Chain{
		chain: chain{
			simplices: []*Simplex{},
		},
		complex:    c.complex,
		chaingroup: c.chaingroup,
		dim:        c.dim,
		idxs:       map[Index]*Simplex{},
		base:       map[Index]struct{}{},
	}

	// only keep the simplices that show up twice
	for smplx, count := range simplexCount {
		if count != 2 {
			continue
		}

		chain.simplices = append(chain.simplices, smplx)
		chain.idxs[smplx.index] = smplx
		for _, v := range smplx.base {
			chain.base[v] = struct{}{}
		}
	}

	if len(chain.simplices) == 0 {
		return c.chaingroup.zero
	}

	return chain
}

// AddSimplex is a convenience method for adding the Chain containing only the Simplex s to the Chain c.
func (c *Chain) AddSimplex(s *Simplex) *Chain {
	return c.Add(c.chaingroup.Singleton(s))
}

// Simplices returns a copy of the simplices that make up c.
func (c *Chain) Simplices() []*Simplex {
	if !c.sorted {
		c.sort()
	}

	simplices := make([]*Simplex, len(c.simplices))
	copy(simplices, c.simplices)

	return simplices
}

// Vector returns the vector representation of c.
// A Chain c is represented as a Vector v by assigning v_i = 1 if
// c contains the i^th simplex in the basis of the ChainGroup; v_i = 0 otherwise.
func (c *Chain) Vector() Vector {
	if c.vector != nil {
		v := mat.DenseCopyOf(c.vector)
		return mat.Matrix(v).(Vector)
	}

	rank := c.chaingroup.Rank()
	c.vector = mat.NewDense(rank, 1, nil)
	n := Index(rank)

	v := c.vector.(*mat.Dense)

	for i := Index(0); i < n; i++ {
		if _, inChain := c.idxs[i]; inChain {
			v.Set(int(i), 0, 1.0)
		}
	}

	vv := mat.DenseCopyOf(c.vector)
	return mat.Matrix(vv).(Vector)
}

// Boundary is a group homomorphism from a p-dimensional ChainGroup to a (p-1)-dimensional ChainGroup.
// In particular, Boundary returns the Chain of simplices that make up the boundary/faces of c.
// For example: If c represents an edge, then the boundary is the chain consisting of the 2 vertices that it connects; if c is a filled in triangle, the boundary is the chain of the 3 edges that make up the triangle.
//
// More info: https://en.wikipedia.org/wiki/Simplicial_homology#Boundaries_and_cycles
func (c *Chain) Boundary() *Chain {
	if c.dim == 0 {
		return nil
	}

	if c.isCycle {
		return c.chaingroup.zero
	}

	if !c.sorted {
		c.sort()
	}

	complex := c.complex
	group := c.chaingroup
	lowerGroup := complex.ChainGroup(group.dim - 1)

	bm := group.BoundaryMap().BoundaryMatrix()
	v := c.Vector().(mat.Matrix)

	bmm, _ := bm.Dims()

	x := mat.NewDense(bmm, 1, nil)
	x.Mul(bm, v)

	for i := 0; i < bmm; i++ {
		x.Set(i, 0, float64(int(x.At(i, 0))%2))
	}

	xv := mat.Matrix(x).(Vector)

	boundary := lowerGroup.ChainFromVector(xv)
	boundary.isCycle = true

	return boundary
}

// Equals returns true is c is equal to a.
func (c *Chain) Equals(a *Chain) bool {
	if c.complex != a.complex || c.dim != a.dim {
		return false
	}

	if len(c.simplices) != len(a.simplices) {
		return false
	}

	c.sort()
	a.sort()

	for idx := range c.simplices {
		cs := c.simplices[idx]
		as := a.simplices[idx]
		if !cs.Equals(as) {
			return false
		}
	}

	return true
}
