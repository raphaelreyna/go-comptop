package comptop

type SimplicialSet struct {
	set map[*Simplex]struct{}

	slice     []*Simplex
	eulerChar *int
}

func NewSimplicialSet(smplxs ...*Simplex) *SimplicialSet {
	ss := &SimplicialSet{
		set: map[*Simplex]struct{}{},
	}

	for _, smplx := range smplxs {
		ss.set[smplx] = struct{}{}
	}

	return ss
}

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

func (ss *SimplicialSet) Ord() int {
	return len(ss.set)
}

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
