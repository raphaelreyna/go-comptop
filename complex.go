package comptop

import "fmt"

type Dim uint

type Complex struct {
	dim Dim

	chainGroups ChainGroups
}

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

func (c *Complex) String() string {
	s := fmt.Sprintf(`Complex{"dim": %d`, c.dim)

	return s
}
