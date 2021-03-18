package comptop

import "fmt"

type ChainGroup struct {
	complex *Complex

	elements map[Index]*Simplex
	dim      Dim

	head Index
}

func (c *Complex) newChainGroup(dim Dim) *ChainGroup {
	return &ChainGroup{
		complex:  c,
		elements: map[Index]*Simplex{},
		dim:      dim,
	}
}

func (cg *ChainGroup) addElement(s *Simplex) {
	if s.Dim() != cg.dim {
		return
	}

	s.index = cg.head
	if _, exists := cg.elements[s.index]; exists {
		return
	}

	cg.head++

	cg.elements[s.index] = s
}

func (cg *ChainGroup) Dim() Dim {
	return cg.dim
}

func (cg *ChainGroup) Ord() int {
	return len(cg.elements)
}

func (cg *ChainGroup) Elements() []*Simplex {
	els := []*Simplex{}

	for _, s := range cg.elements {
		els = append(els, s)
	}

	return els
}

func (cg *ChainGroup) Element(idx Index) *Simplex {
	return cg.elements[idx]
}

func (cg *ChainGroup) String() string {
	return fmt.Sprintf(`ChainGroup{dim: %d, ord: %d, elements: %v}`,
		cg.dim,
		len(cg.elements),
		cg.elements,
	)
}

type ChainGroups map[Dim]*ChainGroup
