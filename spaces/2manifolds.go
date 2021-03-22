package spaces

import comptop "github.com/raphaelreyna/go-comptop"

var Torus []comptop.Base = []comptop.Base{
	{0, 1, 4},
	{1, 4, 5},
	{1, 2, 5},
	{2, 5, 6},
	{0, 2, 6},
	{0, 4, 6},
	{4, 5, 7},
	{5, 7, 8},
	{5, 6, 8},
	{6, 8, 9},
	{4, 6, 9},
	{4, 7, 9},
	{0, 7, 8},
	{0, 1, 8},
	{1, 8, 9},
	{1, 2, 9},
	{2, 7, 9},
	{0, 2, 7},
}

var Anulus []comptop.Base = []comptop.Base{
	{0, 1, 2},
	{1, 2, 3},
	{2, 3, 4},
	{3, 4, 5},
	{0, 4, 5},
	{0, 1, 5},
}

var MobiusStrip []comptop.Base = []comptop.Base{
	{0, 1, 2},
	{1, 2, 3},
	{2, 3, 4},
	{3, 4, 5},
	{1, 4, 5},
	{0, 1, 5},
}
