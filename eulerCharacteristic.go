package comptop

func (s *Simplex) EulerChar() int {
	return 1
}

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

func (ss *SimplicialSet) EulerChar() int {
	if ss.eulerChar != nil {
		return *ss.eulerChar
	}

	var (
		m int = 0
		a int = 1
	)

	rs := ss.RankedSlices()

	for d := Dim(0); d <= Dim(len(rs)); d++ {
		m += a * len(rs[d])
		a *= -1
	}

	ss.eulerChar = &m

	return m
}
