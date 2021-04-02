package spaces

import comptop "github.com/raphaelreyna/go-comptop"

var Torus []comptop.Base = []comptop.Base{
	{0, 1, 3},
	{1, 3, 4},
	{1, 2, 4},
	{2, 4, 5},
	{0, 2, 5},
	{0, 3, 5},

	{3, 4, 6},
	{4, 6, 7},
	{4, 5, 7},
	{5, 7, 8},
	{3, 5, 8},
	{3, 6, 8},

	{0, 6, 7},
	{0, 1, 7},
	{1, 7, 8},
	{1, 2, 8},
	{2, 6, 8},
	{0, 2, 6},
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
