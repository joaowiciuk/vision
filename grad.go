package vision

import (
	"image"
	"math"

	"github.com/joaowiciuk/matrix"
	"github.com/joaowiciuk/vision/kernel"
)

// Grad computes the grad and returns its magnitude and angle.
func Grad(gray *image.Gray) (mag, ang *image.Gray) {
	k := 3
	src := Gray2Mat(gray)
	dx := kernel.SobelX(k)
	dy := kernel.SobelY(k)
	ma, na := dx.Size()
	mb, nb := src.Size()
	cxa, cya := dx.Center()
	top, left, bottom, right := cya, cxa, ma-cya-1, na-cxa-1
	x := func(c, r int) float64 {
		if c >= 0 && c <= nb-1 && r >= 0 && r <= mb-1 {
			return (*src)[r][c] //Inside
		} else if c < 0 && r >= 0 && r <= mb-1 {
			return (*src)[r][0] //Left
		} else if r < 0 && c >= 0 && c <= nb-1 {
			return (*src)[0][c] //Top
		} else if c > nb-1 && r >= 0 && r <= mb-1 {
			return (*src)[r][nb-1] //Right
		} else if r > mb-1 && c >= 0 && c <= nb-1 {
			return (*src)[mb-1][c] //Bottom
		} else if c < 0 && r > mb-1 {
			return (*src)[mb-1][0] //Bottom left corner
		} else if c < 0 && r < 0 {
			return (*src)[0][0] //Top left corner
		} else if c > nb-1 && r < 0 {
			return (*src)[0][nb-1] //Top right corner
		} else {
			return (*src)[mb-1][nb-1] //Bottom right corner
		}
	}
	h := func(c, r int) (h1 float64, h2 float64) {
		m, n := -r+cya, -c+cxa
		if n < 0 || m < 0 {
			return 1, 1
		}
		if n >= na || m >= ma {
			return 1, 1
		}
		return (*dx)[m][n], (*dy)[m][n]
	}
	y := func(c, r int) (y1 float64, y2 float64) {
		r0, r1 := r-top, r+bottom
		c0, c1 := c-left, c+right
		sum1 := 0.
		sum2 := 0.
		for j := r0; j <= r1; j++ {
			for i := c0; i <= c1; i++ {
				h1, h2 := h(c-i, r-j)
				sum1 += x(i, j) * h1
				sum2 += x(i, j) * h2
			}
		}
		return sum1, sum2
	}
	magMat, angMat := matrix.New(mb, nb), matrix.New(mb, nb)
	for r := 0; r < mb; r++ {
		for c := 0; c < nb; c++ {
			h, v := y(c, r)
			(*magMat)[r][c] = math.Hypot(h, v)

			(*angMat)[r][c] = math.Atan2(v, h) * 180 / math.Pi
			for (*angMat)[r][c] < 0 {
				(*angMat)[r][c] += 180
			}
		}
	}
	mag = Mat2Gray(magMat)
	ang = Mat2Gray(angMat)
	return
}
