// The methods in this file satisfy the sort.Interface interface.
// Sorting a simplex sorts its base of indices.
package comptop

func (s *simplex) Len() int {
	return len(s.base)
}

func (s *simplex) Less(i, j int) bool {
	return s.base[i] < s.base[j]
}

func (s *simplex) Swap(i, j int) {
	s.base[i], s.base[j] = s.base[j], s.base[i]
}
