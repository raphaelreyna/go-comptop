package comptop

// EulerChar returns the Euler characteristic of s.
// This function is always equal to 1 for every Simplex.
//
// More info: https://en.wikipedia.org/wiki/Euler_characteristic
func (s *Simplex) EulerChar() int {
	return 1
}

// EulerChar returns the Euler characteristic of c.
// For chains, this value coincides with the number of connected components in the chain.
//
// More info: https://en.wikipedia.org/wiki/Euler_characteristic
func (c *Chain) EulerChar() int {
	if c.eulerChar != nil {
		return *c.eulerChar
	}

	simplices := map[Dim]map[*Simplex]struct{}{}
	dim := c.dim

	stack := make([]*Simplex, len(c.simplices))
	copy(stack, c.simplices)

	// Organize all simplices by dimension
	for len(stack) != 0 {
		n := len(stack) - 1
		smplx := stack[n]
		stack = stack[:n]

		set := simplices[smplx.Dim()]
		if set == nil {
			set = map[*Simplex]struct{}{}
			simplices[smplx.Dim()] = set
		}
		set[smplx] = struct{}{}

		if smplx.Dim() == 0 {
			continue
		}

		faces := smplx.Faces(smplx.dim() - 1)
		for _, face := range faces.Slice() {
			stack = append(stack, face)
		}
	}

	var (
		m int = 0
		a int = 1
	)

	for d := Dim(0); d <= dim; d++ {
		m += a * len(simplices[d])
		a *= -1
	}

	c.eulerChar = &m

	return m
}

// EulerChar returns the Euler characteristic of the Complex.
//
// More info: https://en.wikipedia.org/wiki/Euler_characteristic
func (c *Complex) EulerChar() int {
	if c.eulerChar != nil {
		return *c.eulerChar
	}

	var (
		m int = 0
		a int = 1
	)

	for d := Dim(0); d <= c.dim; d++ {
		m += a * len(c.chainGroups[d].simplices)
		a *= -1
	}

	c.eulerChar = &m

	return m
}

// EulerChar returns the Euler characteristic of the SimplicialSet.
//
// More info: https://en.wikipedia.org/wiki/Euler_characteristic
func (ss *SimplicialSet) EulerChar() int {
	if ss.eulerChar != nil {
		return *ss.eulerChar
	}

	var (
		m int = 0
		a int = 1
	)

	allSimplices := map[Dim]map[*Simplex]struct{}{}
	for smplx := range ss.set {
		if _, exists := allSimplices[smplx.Dim()]; !exists {
			allSimplices[smplx.Dim()] = map[*Simplex]struct{}{}
		}
		allSimplices[smplx.Dim()][smplx] = struct{}{}

		for d := Dim(0); d < smplx.Dim(); d++ {
			faces := smplx.Faces(d)
			for face := range faces.set {
				dd := face.Dim()
				if _, exists := allSimplices[dd]; !exists {
					allSimplices[dd] = map[*Simplex]struct{}{}
				}
				allSimplices[dd][face] = struct{}{}
			}
		}
	}

	for d := Dim(0); d < Dim(len(allSimplices)); d++ {
		m += a * len(allSimplices[d])
		a *= -1
	}

	ss.eulerChar = &m

	return m
}
