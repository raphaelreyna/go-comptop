package comptop

func (c *Complex) LevelSet(f CF, s int) *Complex {
	cmplx := &Complex{}

	for d := Dim(0); d <= c.dim; d++ {
		for _, smplx := range c.GetdSimplices(c.dim - d) {
			if cmplx.GetSimplex(smplx.base...) != nil {
				continue
			}
			inSet := true
			for _, idx := range smplx.Base() {
				if f(idx) != s {
					inSet = false
					break
				}
			}
			if inSet {
				ss := cmplx.NewSimplex(smplx.base...)
				ss.Data = smplx.Data
			}
		}
	}

	return cmplx
}

func (c *Complex) UpperExcursionSet(f CF, s int) *Complex {
	cmplx := &Complex{}

	for d := Dim(0); d <= c.dim; d++ {
		for _, smplx := range c.GetdSimplices(c.dim - d) {
			if cmplx.GetSimplex(smplx.base...) != nil {
				continue
			}
			inSet := true
			for _, idx := range smplx.Base() {
				if f(idx) <= s {
					inSet = false
					break
				}
			}
			if inSet {
				ss := cmplx.NewSimplex(smplx.base...)
				ss.Data = smplx.Data
			}
		}
	}

	return cmplx
}

func (c *Complex) LowerExcursionSet(f CF, s int) *Complex {
	cmplx := &Complex{}

	for d := Dim(0); d <= c.dim; d++ {
		for _, smplx := range c.GetdSimplices(c.dim - d) {
			if cmplx.GetSimplex(smplx.base...) != nil {
				continue
			}
			inSet := true
			for _, idx := range smplx.Base() {
				if f(idx) >= s {
					inSet = false
					break
				}
			}
			if inSet {
				ss := cmplx.NewSimplex(smplx.base...)
				ss.Data = smplx.Data
			}
		}
	}

	return cmplx
}

func (c *Complex) EulerIntegral(a, b int, f CF) int {
	x := 0
	for s := a; s <= b; s++ {
		set := c.UpperExcursionSet(f, s)
		if set == nil {
			continue
		}

		x += set.EulerChar()
	}

	return x
}
