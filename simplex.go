package comptop

import (
	"fmt"
	"sort"
)

type Index uint

type Base []Index

type simplex struct {
	base   Base
	sorted bool
}

func (s *simplex) String() string {
	s.sort()
	return fmt.Sprintf("%+v", s.base)
}

func (s *simplex) sort() {
	if s.sorted {
		return
	}

	sort.Sort(s)
	s.sorted = true
}

func (s *simplex) equals(f *simplex) bool {
	sortsimplices(s, f)

	if len(s.base) != len(f.base) {
		return false
	}

	for idx := range s.base {
		if s.base[idx] != f.base[idx] {
			return false
		}
	}

	return true
}

func (s *simplex) dim() Dim {
	return Dim(len(s.base)) - 1
}

func (s *simplex) d() []*simplex {
	sortsimplices(s)

	boundary := []*simplex{}
	n := len(s.base)

	for j := 0; j < n; j++ {
		f := make([]Index, n-1)
		copy(f[0:j], s.base[0:j])
		copy(f[j:n-1], s.base[j+1:n])

		ss := &simplex{
			base:   f,
			sorted: false,
		}

		boundary = append(boundary, ss)
	}

	return boundary
}

type Simplex struct {
	simplex

	complex *Complex
	index   Index

	faces map[Dim]*SimplicialSet

	Data interface{}
}

func (c *Complex) NewSimplex(base ...Index) *Simplex {
	if c.chainGroups == nil {
		c.chainGroups = ChainGroups{}
	}

	dim := Dim(len(base)) - 1

	if dim > c.dim {
		c.dim = dim
	}

	s := &simplex{base: base}
	var newSimplex *Simplex

	stack := []*simplex{s}

	for len(stack) > 0 {
		n := len(stack) - 1

		// pop next simplex from stack
		ss := stack[n]
		stack = stack[:n]

		// Skip this simplex if its already in the complex
		if smplx := c.GetSimplex(ss.base...); smplx != nil {
			if smplx.Dim() == dim {
				return smplx
			}
			continue
		}

		// Add this simplex to the appropriate chain group
		p := ss.dim()
		smplx := &Simplex{
			simplex: *ss,
			complex: c,
		}

		group := c.chainGroups[p]
		if group == nil {
			group = c.newChainGroup(p)
			c.chainGroups[p] = group
		}
		group.addSimplex(smplx)

		if newSimplex == nil {
			newSimplex = smplx
		}

		if p == 0 {
			continue
		}

		// Compute the boundary and add all its simplices to the stack
		for _, sss := range ss.d() {
			stack = append(stack, sss)
		}
	}

	// cached results should be reset
	c.eulerChar = nil
	c.strng = ""

	return newSimplex
}

func (s *Simplex) String() string {
	return fmt.Sprintf(
		`Simplex{"dim": %d, "index": %d, "base": %v, "data": %+v}`,
		s.Dim(),
		s.index,
		s.base,
		s.Data,
	)
}

func (s *Simplex) Dim() Dim {
	return Dim(len(s.base)) - 1
}

func (s *Simplex) Index() Index {
	return s.index
}

func (s *Simplex) Equals(f *Simplex) bool {
	if s.complex != f.complex {
		return false
	}

	return s.simplex.equals(&f.simplex)
}

func (s *Simplex) Sort() {
	s.simplex.sort()
}

func (s *Simplex) HasFace(f *Simplex) bool {
	if f.Dim() != s.Dim()-1 {
		return false
	}

	boundary := s.d()

	for _, smplx := range boundary {
		ss := &Simplex{simplex: *smplx, complex: s.complex}
		if ss.Equals(f) {
			return true
		}
	}

	return false
}

func (s *Simplex) Intersection(g *Simplex) *Simplex {
	sortSimplices(s, g)

	intersection := make([]Index, 0)

	n := g.Len()

	for _, el := range s.base {
		idx := sort.Search(n, func(j int) bool {
			return g.base[j] >= el
		})
		if idx < n && g.base[idx] == el {
			intersection = append(intersection, el)
		}
	}

	return s.complex.GetSimplex(intersection...)
}

func (s *Simplex) Boundary() *Chain {
	chain := &Chain{simplices: []*Simplex{}, dim: s.Dim() - 1}

	for _, smplx := range s.d() {
		ss := s.complex.GetSimplex(smplx.base...)
		chain.simplices = append(chain.simplices, ss)
	}

	return chain
}

func (s *Simplex) Faces(d Dim) *SimplicialSet {
	if faces, exists := s.faces[d]; exists {
		return faces
	}

	faces := []*Simplex{}
	complex := s.complex
	group := complex.chainGroups[d]

	for _, smplx := range group.simplices {
		if s.HasFace(smplx) {
			faces = append(faces, smplx)
		}
	}

	if s.faces == nil {
		s.faces = map[Dim]*SimplicialSet{}
	}

	s.faces[d] = NewSimplicialSet(faces...)

	return s.faces[d]
}

func sortsimplices(simplices ...*simplex) {
	for _, s := range simplices {
		s.sort()
	}
}

func sortSimplices(simplices ...*Simplex) {
	for _, s := range simplices {
		s.Sort()
	}
}
