package comptop

import (
	"fmt"
)

// Complex represents an abstract simplicial complex.
// More info: https://en.wikipedia.org/wiki/Abstract_simplicial_complex
type Complex struct {
	dim Dim

	chainGroups ChainGroups
	principles  map[*Simplex]struct{}

	eulerChar *int

	strng string
}

func (c *Complex) chaingroup(d Dim) *ChainGroup {
	if d > c.dim {
		return nil
	}

	if c.chainGroups == nil {
		c.chainGroups = ChainGroups{}
	}

	if c.chainGroups[d] == nil {
		c.chainGroups[d] = c.newChainGroup(d)
	}

	return c.chainGroups[d]
}

// GetSimplex returns the Simplex consisting of 0-simplices with base indices.
func (c *Complex) GetSimplex(base ...Index) *Simplex {
	s := &Simplex{complex: c}
	s.base = base

	dim := s.Dim()

	if c.chainGroups == nil {
		c.chainGroups = ChainGroups{}
		return nil
	}

	var group *ChainGroup
	if group = c.chainGroups[dim]; group == nil {
		return nil
	}

	for _, smplx := range group.simplices {
		if smplx.Equals(s) {
			return smplx
		}
	}

	return nil
}

// GetSimplexByIndex returns the Simplex of dimension d with Index idx.
func (c *Complex) GetSimplexByIndex(idx Index, d Dim) *Simplex {
	group := c.chainGroups[d]
	if group == nil {
		return nil
	}

	return group.simplices[idx]
}

func (c *Complex) GetdSimplices(d Dim) []*Simplex {
	if d > c.dim {
		return nil
	}

	g := c.chaingroup(d)

	return g.Simplices()
}

// ChainGroup returns the free abelian group of d-chains in the Complex.
//
// More info: https://en.wikipedia.org/wiki/Free_abelian_group
func (c *Complex) GetChainGroup(d Dim) *ChainGroup {
	return c.chaingroup(d)
}

func (c *Complex) String() string {
	if c.strng != "" {
		return c.strng
	}

	s := fmt.Sprintf(`Complex{"dim": %d,`, c.dim)
	s += ` "simplices": {`

	if c.principles == nil {
		c.principles = c.PrincipleSimplices().set
	}

	for smplx := range c.principles {
		s += smplx.String() + ", "
	}

	s += "}}"

	c.strng = s

	return s
}

// PrincipleSimplices returns the set of principle simplices in the Complex.
// A Simplex is principle if it's not the face of any other Simplex (has no cofaces).
func (c *Complex) PrincipleSimplices() *SimplicialSet {
	if c.principles != nil {
		return &SimplicialSet{
			set: c.principles,
		}
	}

	p := map[*Simplex]struct{}{}

	for d := int(c.dim); d >= 0; d-- {
		group := c.chainGroups[Dim(d)]
		for _, smplx := range group.simplices {
			if d == int(c.dim) {
				p[smplx] = struct{}{}
				continue
			}

			higherGroup := c.chainGroups[Dim(d)+1]
			var isFace bool
			for _, hsmplx := range higherGroup.simplices {
				if isFace = hsmplx.HasFace(smplx); isFace {
					break
				}
			}

			if !isFace {
				p[smplx] = struct{}{}
			}
		}
	}

	c.principles = p

	return &SimplicialSet{set: p}
}

func (c *Complex) ReducedBettiNumbers() []int {
	var (
		z     int
		betti = []int{}
	)

	for dim := Dim(0); dim <= c.dim; dim++ {
		g := c.GetChainGroup(dim)
		if z != 0 {
			betti = append(betti, z-g.BoundaryMap().SmithNormalDiagonalLength())
		}
		bm := c.GetChainGroup(dim).BoundaryMap()
		z = bm.Zp()
	}

	return append(betti, z)
}

func (c *Complex) BettiNumbers() []int {
	rb := c.ReducedBettiNumbers()
	rb[0]++
	return rb
}

func (c *Complex) resetCache() {
	c.eulerChar = nil
	c.strng = ""
	c.principles = nil
}
