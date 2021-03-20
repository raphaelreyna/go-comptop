package comptop

// SimplicialSet is a set containing simplices.
// Note: SimplicialSet is not a category theoretical simplicial set.
type SimplicialSet struct {
	set map[*Simplex]struct{}

	slice     []*Simplex
	eulerChar *int
}

// NewSimplicialSet returns the SimplicialSet containing the provided simplices.
func NewSimplicialSet(smplxs ...*Simplex) *SimplicialSet {
	ss := &SimplicialSet{
		set: map[*Simplex]struct{}{},
	}

	for _, smplx := range smplxs {
		ss.set[smplx] = struct{}{}
	}

	return ss
}

// Add appends the provided simplices to the SimplicialSet.
func (ss *SimplicialSet) Add(smplxs ...*Simplex) {
	if ss.set == nil {
		ss.set = map[*Simplex]struct{}{}
	}

	for _, smplx := range smplxs {
		ss.set[smplx] = struct{}{}
	}

	ss.slice = nil
	ss.eulerChar = nil
}

// Rem removes the provided elements from the SimplicialSet.
func (ss *SimplicialSet) Rem(smplxs ...*Simplex) {
	if ss.set == nil {
		return
	}

	for _, smplx := range smplxs {
		delete(ss.set, smplx)
	}

	ss.slice = nil
	ss.eulerChar = nil
}

// Slice returns the simplices in the set as a []*Simplex slice.
func (ss *SimplicialSet) Slice() []*Simplex {
	if ss.slice != nil {
		slice := make([]*Simplex, len(ss.slice))
		copy(slice, ss.slice)
		return slice
	}

	ss.slice = []*Simplex{}
	for smplx := range ss.set {
		ss.slice = append(ss.slice, smplx)
	}

	slice := make([]*Simplex, len(ss.slice))
	copy(slice, ss.slice)

	return slice
}

// RankedSlices returns the slices in the set organized by their dimension.
func (ss *SimplicialSet) RankedSlices() map[Dim][]*Simplex {
	m := map[Dim][]*Simplex{}

	for smplx := range ss.set {
		d := smplx.Dim()
		if _, exists := m[d]; !exists {
			m[d] = []*Simplex{}
		}

		m[d] = append(m[d], smplx)
	}

	return m
}

// Card returns the cardinality (number of elements in the set) of the SimplicialSet.
func (ss *SimplicialSet) Card() int {
	return len(ss.set)
}

// Union returns the union of ss with the provided sets.
func (ss *SimplicialSet) Union(sets ...*SimplicialSet) *SimplicialSet {
	s := &SimplicialSet{
		set: map[*Simplex]struct{}{},
	}

	for _, set := range sets {
		for el, v := range set.set {
			s.set[el] = v
		}
	}

	return s
}
