package comptop

import (
	"sort"

	"gonum.org/v1/gonum/mat"
)

type Chain struct {
	complex    *Complex
	chaingroup *ChainGroup

	simplices []*Simplex
	vector    Vector
	base      map[Index]struct{}
	idxs      map[Index]*Simplex
	dim       Dim

	eulerChar *int

	isCycle bool

	sorted bool
}

// sort.Interface interface methods
func (c *Chain) Len() int {
	return len(c.simplices)
}

func (c *Chain) Less(i, j int) bool {
	return c.simplices[i].index < c.simplices[j].index
}

func (c *Chain) Swap(i, j int) {
	c.simplices[i], c.simplices[j] = c.simplices[j], c.simplices[i]
}

func (c *Chain) Sort() {
	sort.Sort(c)
}

func (c *Chain) String() string {
	s := "Chain{"

	for _, smplx := range c.simplices {
		s += smplx.String() + ", "

	}

	s += "}"

	return s
}

func (c *Chain) Dim() Dim {
	return c.dim
}

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
		complex:    c.complex,
		chaingroup: c.chaingroup,
		dim:        c.dim,
		simplices:  []*Simplex{},
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

func (c *Chain) AddSimplex(s *Simplex) *Chain {
	return c.Add(c.chaingroup.Singleton(s))
}

func (c *Chain) Simplices() []*Simplex {
	if !c.sorted {
		c.Sort()
	}

	simplices := make([]*Simplex, len(c.simplices))
	copy(simplices, c.simplices)

	return simplices
}

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

func (c *Chain) Boundary() *Chain {
	if c.isCycle {
		return c.chaingroup.zero
	}

	complex := c.complex
	group := c.chaingroup
	lowerGroup := complex.chainGroups[group.dim-1]
	if lowerGroup == nil {
		complex.chainGroups[group.dim-1] = complex.newChainGroup(group.dim - 1)
		lowerGroup = complex.chainGroups[group.dim-1]
	}

	bm := group.boundaryMatrix
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
