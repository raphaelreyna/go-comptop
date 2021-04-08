package comptop

// NewSimplex adds a Simplex to c.
// All lower dimensional faces of the new Simplex are computed and automatically added to c.
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
		if p == 0 {
			smplx.index = ss.base[0]
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
	c.resetCache()

	return newSimplex
}

// NewSimplex adds multiple simplices to c.
// All lower dimensional faces of each new Simplex are computed and automatically added to c.
func (c *Complex) NewSimplices(bases ...Base) *SimplicialSet {
	if c.chainGroups == nil {
		c.chainGroups = ChainGroups{}
	}

	set := map[*Simplex]struct{}{}

	for _, base := range bases {
		dim := Dim(len(base)) - 1

		if dim > c.dim {
			c.dim = dim
		}

		s := &simplex{base: base}
		var newSimplex *Simplex

		stack := []*simplex{s}
		topLevel := true

		for len(stack) > 0 {
			n := len(stack) - 1

			// pop next simplex from stack
			ss := stack[n]
			stack = stack[:n]

			// Skip this simplex if its already in the complex
			if smplx := c.GetSimplex(ss.base...); smplx != nil {
				if topLevel {
					set[smplx] = struct{}{}
				}
				topLevel = false
				continue
			}
			topLevel = false

			// Add this simplex to the appropriate chain group
			p := ss.dim()
			smplx := &Simplex{
				simplex: *ss,
				complex: c,
			}
			if p == 0 {
				smplx.index = ss.base[0]
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

		if newSimplex != nil {
			set[newSimplex] = struct{}{}
		}
	}

	// cached results should be reset
	c.resetCache()

	return &SimplicialSet{set: set}
}

// DataProvider is used to attach user-defined data to simplices.
type DataProvider func(Dim, Index, Base) interface{}

// NewSimplex adds a Simplex to c while using dp to attach data to each newly created Simplex.
// All lower dimensional faces of the new Simplex are computed and automatically added to c.
func (c *Complex) NewSimplexWithData(dp DataProvider, base ...Index) *Simplex {
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
		if p == 0 {
			smplx.index = ss.base[0]
		}

		group := c.chainGroups[p]
		if group == nil {
			group = c.newChainGroup(p)
			c.chainGroups[p] = group
		}
		group.addSimplex(smplx)

		// Add the data provded by the DataProvider
		smplx.Data = dp(p, smplx.index, base)

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
	c.resetCache()

	return newSimplex
}

// NewSimplex adds multiple simplices to c, using dp to attach data to each newly created Simplex.
// All lower dimensional faces of each new Simplex are computed and automatically added to c.
func (c *Complex) NewSimplicesWithData(dp DataProvider, bases ...Base) *SimplicialSet {
	if c.chainGroups == nil {
		c.chainGroups = ChainGroups{}
	}

	set := map[*Simplex]struct{}{}

	for _, base := range bases {
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
				continue
			}

			// Add this simplex to the appropriate chain group
			p := ss.dim()
			smplx := &Simplex{
				simplex: *ss,
				complex: c,
			}
			if p == 0 {
				smplx.index = ss.base[0]
			}

			group := c.chainGroups[p]
			if group == nil {
				group = c.newChainGroup(p)
				c.chainGroups[p] = group
			}
			group.addSimplex(smplx)

			// Add the data provded by the DataProvider
			smplx.Data = dp(p, smplx.index, base)

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

		set[newSimplex] = struct{}{}
	}

	// cached results should be reset
	c.resetCache()

	return &SimplicialSet{set: set}
}
