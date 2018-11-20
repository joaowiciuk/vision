/*Package vision provides common algorithms used in computer vision*/
package vision

import (
	"image"
	"image/draw"
	"math"

	"github.com/joaowiciuk/vision/kernel"

	"github.com/joaowiciuk/matrix"
)

// Canny implements the popular canny edge detector
func Canny(img image.Image, upperThreshold, lowerThreshold uint8, k int, σ float64) (j *image.Gray) {
	lowT := float64(lowerThreshold)
	uppT := float64(upperThreshold)
	//Convert the input image to a single src matrix
	tensor := Im2Mat(img)
	var src *matrix.Matrix
	switch len(tensor) {
	case 1:
		src = Im2Mat(img)[0]
	case 3, 4:
		src = matrix.New(img.Bounds().Dy(), img.Bounds().Dx())
		src.Law(func(r, c int) float64 {
			return 0.299*(*tensor[0])[r][c] + 0.587*(*tensor[1])[r][c] + 0.114*(*tensor[2])[r][c]
		})
	default:
		return nil
	}

	m, n := src.Size()

	//Output matrix
	out := matrix.New(m, n)

	//Preprocessing to obtain magnitude and angle from source matrix
	ang := matrix.New(m, n)
	mag := matrix.New(m, n)
	preProc(m, n, mag, ang, src, k, σ)
	/* mag0, ang0 := grad(src, k)
	fmt.Println(mag.Dist(mag0, matrix.Norm1))
	fmt.Println(ang.Dist(ang0, matrix.Norm1)) */

	//Non-maximum suppression
	nonMaxSup(m, n, out, mag, ang, uppT)

	//Hysterysis threshold
	hystThresh(m, n, out, mag, ang, lowT)

	aux := *Mat2Im([]*matrix.Matrix{out})
	j = image.NewGray(aux.Bounds())
	draw.Draw(j, j.Bounds(), aux, j.Bounds().Min, draw.Src)
	return
}

func sobel(k int, B *matrix.Matrix) (magX, magY *matrix.Matrix) {
	dx := kernel.SobelX(k)
	dy := kernel.SobelY(k)

	ma, na := dx.Size()
	mb, nb := B.Size()
	cxa, cya := dx.Center()
	top, left, bottom, right := cya, cxa, ma-cya-1, na-cxa-1
	x := func(c, r int) float64 {
		if c >= 0 && c <= nb-1 && r >= 0 && r <= mb-1 {
			return (*B)[r][c] //Inside
		} else if c < 0 && r >= 0 && r <= mb-1 {
			return (*B)[r][0] //Left
		} else if r < 0 && c >= 0 && c <= nb-1 {
			return (*B)[0][c] //Top
		} else if c > nb-1 && r >= 0 && r <= mb-1 {
			return (*B)[r][nb-1] //Right
		} else if r > mb-1 && c >= 0 && c <= nb-1 {
			return (*B)[mb-1][c] //Bottom
		} else if c < 0 && r > mb-1 {
			return (*B)[mb-1][0] //Bottom left corner
		} else if c < 0 && r < 0 {
			return (*B)[0][0] //Top left corner
		} else if c > nb-1 && r < 0 {
			return (*B)[0][nb-1] //Top right corner
		} else {
			return (*B)[mb-1][nb-1] //Bottom right corner
		}
	}
	h := func(c, r int) (h1 float64, h2 float64) {
		m, n := -r+cya, -c+cxa
		if n < 0 || m < 0 {
			/* fmt.Printf("h[%d, %d] = %.2f\n", c, r, 0.) */
			return 1, 1
		}
		if n >= na || m >= ma {
			/* fmt.Printf("h[%d, %d] = %.2f\n", c, r, 0.) */
			return 1, 1
		}
		/* fmt.Printf("h[%d, %d] = %.2f\n", c, r, (*A)[m][n]) */
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
	magX, magY = matrix.New(mb, nb), matrix.New(mb, nb)
	for r := 0; r < mb; r++ {
		for c := 0; c < nb; c++ {
			(*magX)[r][c], (*magY)[r][c] = y(c, r)
		}
	}
	return
}

func sobelFast(k int, B *matrix.Matrix, c chan func() (*matrix.Matrix, *matrix.Matrix)) {
	c <- (func() (*matrix.Matrix, *matrix.Matrix) {
		dx, dy := sobel(k, B)
		return dx, dy
	})
}

func grad(src *matrix.Matrix, k int) (mag, ang *matrix.Matrix) {
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
			/* fmt.Printf("h[%d, %d] = %.2f\n", c, r, 0.) */
			return 1, 1
		}
		if n >= na || m >= ma {
			/* fmt.Printf("h[%d, %d] = %.2f\n", c, r, 0.) */
			return 1, 1
		}
		/* fmt.Printf("h[%d, %d] = %.2f\n", c, r, (*A)[m][n]) */
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

	mag, ang = matrix.New(mb, nb), matrix.New(mb, nb)
	for r := 0; r < mb; r++ {
		for c := 0; c < nb; c++ {
			h, v := y(c, r)
			(*mag)[r][c] = math.Hypot(h, v)

			(*ang)[r][c] = math.Atan2(v, h) * 180 / math.Pi
			for (*ang)[r][c] < 0 {
				(*ang)[r][c] += 180
			}
		}
	}
	return
}

func nonMaxSup(m, n int, out, mag, ang *matrix.Matrix, uppT float64) {
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			ang0 := (*ang)[y][x]
			if (*mag)[y][x] < uppT {
				continue
			}
			flag := true

			if ang0 > 112.5 && ang0 <= 157.5 {
				if y > 0 && x < n-1 && (*mag)[y][x] <= (*mag)[y-1][x+1] {
					flag = false
				}
				if y < m-1 && x > 0 && (*mag)[y][x] <= (*mag)[y+1][x-1] {
					flag = false
				}
			} else if ang0 > 67.5 && ang0 <= 112.5 {
				if y > 0 && (*mag)[y][x] <= (*mag)[y-1][x] {
					flag = false
				}
				if y < m-1 && (*mag)[y][x] <= (*mag)[y+1][x] {
					flag = false
				}
			} else if ang0 > 22.5 && ang0 <= 67.5 {
				if y > 0 && x > 0 && (*mag)[y][x] <= (*mag)[y-1][x-1] {
					flag = false
				}
				if y < m-1 && x < n-1 && (*mag)[y][x] <= (*mag)[y+1][x+1] {
					flag = false
				}
			} else {
				if x > 0 && (*mag)[y][x] <= (*mag)[y][x-1] {
					flag = false
				}
				if x < n-1 && (*mag)[y][x] <= (*mag)[y][x+1] {
					flag = false
				}
			}

			if flag {
				(*out)[y][x] = 255.
			}
		}
	}
}

func hystThresh(m, n int, out, mag, ang *matrix.Matrix, lowT float64) {
	imageChanged := true
	i := 0
	for imageChanged {
		imageChanged = false
		i++
		for x := 0; x < n; x++ {
			for y := 0; y < m; y++ {
				if x < 2 || x >= n-2 || y < 2 || y >= m-2 {
					continue
				}
				ang0 := (*ang)[y][x]
				if (*out)[y][x] == 255. {
					(*out)[y][x] = 64.
					if ang0 > 112.5 && ang0 <= 157.5 {

						if y > 0 && x > 0 {
							if lowT <= (*mag)[y-1][x-1] &&
								(*out)[y-1][x-1] != 64. &&
								(*ang)[y-1][x-1] > 112.5 &&
								(*ang)[y-1][x-1] <= 157.5 &&
								(*mag)[y-1][x-1] > (*mag)[y-2][x] &&
								(*mag)[y-1][x-1] > (*mag)[y][x-2] {

								(*out)[y-1][x-1] = 255.
								imageChanged = true
							}
						}

						if y < m-1 && x < n-1 {
							if lowT <= (*mag)[y+1][x+1] &&
								(*out)[y+1][x+1] != 64. &&
								(*ang)[y+1][x+1] > 112.5 &&
								(*ang)[y+1][x+1] <= 157.5 &&
								(*mag)[y+1][x+1] > (*mag)[y+2][x] &&
								(*mag)[y+1][x+1] > (*mag)[y][x+2] {

								(*out)[y+1][x+1] = 255.
								imageChanged = true
							}
						}

					} else if ang0 > 67.5 && ang0 <= 112.5 {

						if x > 0 {
							if lowT <= (*mag)[y][x-1] &&
								(*out)[y][x-1] != 64. &&
								(*ang)[y][x-1] > 67.5 &&
								(*ang)[y][x-1] <= 112.5 &&
								(*mag)[y][x-1] > (*mag)[y-1][x-1] &&
								(*mag)[y][x-1] > (*mag)[y+1][x-1] {

								(*out)[y][x-1] = 255.
								imageChanged = true
							}
						}

						if x < n-1 {
							if lowT <= (*mag)[y][x+1] &&
								(*out)[y][x+1] != 64. &&
								(*ang)[y][x+1] > 67.5 &&
								(*ang)[y][x+1] <= 112.5 &&
								(*mag)[y][x+1] > (*mag)[y-1][x+1] &&
								(*mag)[y][x+1] > (*mag)[y+1][x+1] {

								(*out)[y][x+1] = 255.
								imageChanged = true
							}
						}

					} else if ang0 > 22.5 && ang0 <= 67.5 {

						if y > 0 && x < n-1 {
							if lowT <= (*mag)[y-1][x+1] &&
								(*out)[y-1][x+1] != 64. &&
								(*ang)[y-1][x+1] > 22.5 &&
								(*ang)[y-1][x+1] <= 67.5 &&
								(*mag)[y-1][x+1] > (*mag)[y-2][x] &&
								(*mag)[y-1][x+1] > (*mag)[y][x+2] {

								(*out)[y-1][x+1] = 255.
								imageChanged = true
							}
						}

						if y < m-1 && x > 0 {
							if lowT <= (*mag)[y+1][x-1] &&
								(*out)[y+1][x-1] != 64. &&
								(*ang)[y+1][x-1] > 22.5 &&
								(*ang)[y+1][x-1] <= 67.5 &&
								(*mag)[y+1][x-1] > (*mag)[y][x-2] &&
								(*mag)[y+1][x-1] > (*mag)[y+2][x] {

								(*out)[y+1][x-1] = 255.
								imageChanged = true
							}
						}

					} else {
						if y > 0 {
							if lowT <= (*mag)[y-1][x] &&
								(*out)[y-1][x] != 64. &&
								(*ang)[y-1][x] < 22.5 &&
								(*ang)[y-1][x] >= 157.5 &&
								(*mag)[y-1][x] > (*mag)[y-1][x-1] &&
								(*mag)[y-1][x] > (*mag)[y-1][x+2] {

								(*out)[y-1][x] = 255.
								imageChanged = true
							}
						}

						if y < m-1 {
							if lowT <= (*mag)[y+1][x] &&
								(*out)[y+1][x] != 64. &&
								(*ang)[y+1][x] < 22.5 &&
								(*ang)[y+1][x] >= 157.5 &&
								(*mag)[y+1][x] > (*mag)[y+1][x-1] &&
								(*mag)[y+1][x] > (*mag)[y+1][x+1] {

								(*out)[y+1][x] = 255.
								imageChanged = true
							}
						}
					}
				}
			}
		}
	}

	//Reassign
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			if (*out)[y][x] == 64. {
				(*out)[y][x] = 255.
			}
		}
	}
}

func preProc(m, n int, mag, ang, src *matrix.Matrix, k int, σ float64) {
	srcFuture := kernel.Gaussian(k, σ).ConvFast(src)
	/* srcFuture := make(chan *matrix.Matrix, 1)
	srcFuture <- src */
	c := make(chan func() (*matrix.Matrix, *matrix.Matrix))
	go sobelFast(k, <-srcFuture, c)

	magX, magY := (<-c)()

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			(*mag)[r][c] = math.Hypot((*magX)[r][c], (*magY)[r][c])

			(*ang)[r][c] = math.Atan2((*magY)[r][c], (*magX)[r][c]) * 180 / math.Pi
			for (*ang)[r][c] < 0 {
				(*ang)[r][c] += 180
			}
		}
	}
	return
}
