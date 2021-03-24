package comptop

import (
	"fmt"
)

type Dim uint

// Complex represents an abstract simplicial complex.
// More info: https://en.wikipedia.org/wiki/Abstract_simplicial_complex
type Complex struct {
	dim Dim

	chainGroups ChainGroups

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

func (c *Complex) GetChainGroup(d Dim) *ChainGroup {
	return c.chaingroup(d)
}

func (c *Complex) String() string {
	if c.strng != "" {
		return c.strng
	}

	s := fmt.Sprintf(`Complex{"dim": %d`, c.dim)
	s += ` "simplices": {`

	for d := int(c.dim); d >= 0; d-- {
		group := c.chainGroups[Dim(d)]
		for _, smplx := range group.simplices {
			if d == int(c.dim) {
				s += smplx.String() + ", "
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
				s += smplx.String() + ", "
			}
		}
	}

	s += "}}"

	c.strng = s

	return s
}
