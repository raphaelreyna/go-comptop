package comptop

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// ComplexFromOBJFile reads in the .obj file at path and builds a *Complex
// out of the face data in the file.
func ComplexFromOBJFile(path string) (*Complex, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fr := bufio.NewReader(file)

	var (
		c *Complex = &Complex{}
		b []Base   = []Base{}
		l []byte
	)

	l, _, err = fr.ReadLine()
	for err == nil {
		// Skip empty lines
		if len(l) == 0 {
			continue
		}

		// We're only interested in faces
		if l[0] != 'f' {
			continue
		}

		// Grab the 3 indices that make up a face 2-simplex
		parts := strings.Split(string(l), " ")
		parts = append([]string(nil), parts[1:]...)

		bb := Base{}
		for _, p := range parts {
			// Remove vertex normal and text-coord data
			if p = strings.Split(p, "/")[0]; p == "" {
				continue
			}

			idx, er := strconv.Atoi(p)
			if er != nil {
				return nil, er
			}

			bb = append(bb, Index(idx))
		}

		l, _, err = fr.ReadLine()
	}

	switch err {
	case io.EOF:
		break
	default:
		return nil, err
	}

	c.NewSimplices(b...)

	return c, nil
}
