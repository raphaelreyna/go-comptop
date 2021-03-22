package comptop

import (
	"fmt"
	"sort"
)

// Index is used to establiish a total order amongst simplices of the same dimension.
// Index is also used to uniquely identify simplicies up to dimension.
type Index uint

// Base is a collection of indices for 0-dimensional simplices.
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

// Simplex is a p-dimensional polytope which is the convex hull of its p+1 0-dimensional simplices (points/vertices).
// Every Simplex should be part of a Complex; every Simplex in a Complex is considered to live in the same topological space.
// Simplex is uniquely identified in its Complex by its dimension ((*Simplex).Dim) and Index ((*Simplex).Index).
// Simplex can encapsulate user-defined data in its Data field.
//
// More info: https://encyclopediaofmath.org/wiki/Simplex_(abstract)
type Simplex struct {
	simplex

	complex *Complex
	index   Index

	faces map[Dim]*SimplicialSet

	Data interface{}
}

func (s *Simplex) String() string {
	if s.Data == nil {
		return fmt.Sprintf(
			`Simplex{"dim": %d, "index": %d, "base": %v}`,
			s.Dim(),
			s.index,
			s.base,
		)
	}

	return fmt.Sprintf(
		`Simplex{"dim": %d, "index": %d, "base": %v, "data": %+v}`,
		s.Dim(),
		s.index,
		s.base,
		s.Data,
	)
}

func (s *Simplex) Complex() *Complex {
	return s.complex
}

// Dim returns the dimension of s, which is defined to be 1 + (# of points/0-simplices in s).
func (s *Simplex) Dim() Dim {
	return Dim(len(s.base)) - 1
}

// Index returns the Index of s, which uniquely identifies it in the basis of its corresponding ChainGroup.
func (s *Simplex) Index() Index {
	return s.index
}

// Equal returns true if s and f are equal; returns false otherwise.
func (s *Simplex) Equals(f *Simplex) bool {
	if s == nil || f == nil {
		return false
	}

	if s.complex != f.complex {
		return false
	}

	return s.simplex.equals(&f.simplex)
}

// Sort sorts the base of s by increasing order.
func (s *Simplex) Sort() {
	s.simplex.sort()
}

// Base returns a copy of the base set of s.
func (s *Simplex) Base() Base {
	b := make(Base, len(s.base))
	copy(b, s.base)

	return b
}

// HasFace returns true if s has f as a face.
func (s *Simplex) HasFace(f *Simplex) bool {
	if s.Intersection(f).Equals(f) {
		return true
	}
	return false
}

// Intersection returns the intersection of simplices s and g.
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

// Boundary computes the boundary of s as a Chain in a ChainGroup of the Complex of s.
func (s *Simplex) Boundary() *Chain {
	chain := &Chain{chain: chain{simplices: []*Simplex{}}, dim: s.Dim() - 1}

	for _, smplx := range s.d() {
		ss := s.complex.GetSimplex(smplx.base...)
		chain.simplices = append(chain.simplices, ss)
	}

	return chain
}

// Faces returns the set of d dimensional faces of s.
func (s *Simplex) Faces(d Dim) *SimplicialSet {
	if s.faces == nil {
		s.faces = map[Dim]*SimplicialSet{}
	}

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
