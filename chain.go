package comptop

import (
	"sort"
)

type Chain struct {
	simplices []*Simplex
	dim       Dim
	sorted    bool
}

// sort.Interface interface methods
func (c *Chain) Len() int {
	return len(c.simplices)
}

func (c *Chain) Less(i, j int) bool {
	return c.simplices[i].index < c.simplices[j].index
}

func (c *Chain) Swap(i, j int) {
	c.simplices[i], c.simplices[j] = c.simplices[j], c.simplices[i]
}

func (c *Chain) Sort() {
	sort.Sort(c)
}

func (c *Chain) String() string {
	s := "Chain{"

	for _, smplx := range c.simplices {
		s += smplx.String() + ", "

	}

	s += "}"

	return s
}
