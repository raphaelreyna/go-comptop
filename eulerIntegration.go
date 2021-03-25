package comptop

// LevelSet returns a sub-Complex of c;
// if smplx is a Simplex in Complex c, then smplx is in the level set if f(v) = s for each vertex v in smplx.
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

// UpperExcursionSet returns a sub-Complex of c;
// if smplx is a Simplex in Complex c, then smplx is in the upper excursion set if f(v) > s for each vertex v in smplx.
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

// LowerExcursionSet returns a sub-Complex of c;
// if smplx is a Simplex in Complex c, then smplx is in the upper excursion set if f(v) < s for each vertex v in smplx.
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

// EulerIntegral computes the integral of f for s-values ranging from a to b over the Complex, using the Euler characteristic as its measure.
// For better stability, this method uses upper excursion sets rather than level sets.
//
// We compute the Euler integral as:
//    /-\		      _b_                  _b_
//    |				  \                    \
//    | f(x) d'\/  =  /__ s * '\/({f=s}) = /__ '\/({f>s})
//	  |		   /\,	  s=a      /\,         s=a  /\,
//  \-/ c
//
// where {f=s} is a level set and {f>s} is an upper excursion set.
//
// More info: https://en.wikipedia.org/wiki/Euler_calculus
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
