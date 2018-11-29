package vision

import (
	"image"
	"math"
	"sync"
)

func EdgeDrawing(gray *image.Gray) *image.Gray {
	dx := [][]float64{{1, 0, -1}, {2, 0, -2}, {1, 0, -1}}
	dy := [][]float64{{1, 2, 1}, {0, 0, 0}, {-1, -2, -1}}
	mb, nb := gray.Bounds().Dy(), gray.Bounds().Dx()

	//Extend image signal at borders
	signal := func(x, y int) float64 {
		if x >= 0 && x <= nb-1 && y >= 0 && y <= mb-1 {
			return float64(gray.Pix[y*nb+x]) //Inside
		} else if x < 0 && y >= 0 && y <= mb-1 {
			return float64(gray.Pix[y*nb]) //Left
		} else if y < 0 && x >= 0 && x <= nb-1 {
			return float64(gray.Pix[x]) //Top
		} else if x > nb-1 && y >= 0 && y <= mb-1 {
			return float64(gray.Pix[y*nb+(nb-1)]) //Right
		} else if y > mb-1 && x >= 0 && x <= nb-1 {
			return float64(gray.Pix[(mb-1)*nb+x]) //Bottom
		} else if x < 0 && y > mb-1 {
			return float64(gray.Pix[(mb-1)*nb]) //Bottom left corner
		} else if x < 0 && y < 0 {
			return float64(gray.Pix[0]) //Top left corner
		} else if x > nb-1 && y < 0 {
			return float64(gray.Pix[nb-1]) //Top right corner
		} else {
			return float64(gray.Pix[(mb-1)*nb+(nb-1)]) //Bottom right corner
		}
	}
	sobel := func(x, y int) (sobelX float64, sobelY float64) {
		m, n := -y+1, -x+1
		if n < 0 || m < 0 {
			return 1, 1
		}
		if n >= 3 || m >= 3 {
			return 1, 1
		}
		return dx[m][n], dy[m][n]
	}
	conv := func(x, y int) (convX float64, convY float64) {
		y0, y1 := y-1, y+1
		x0, x1 := x-1, x+1
		convSumX := 0.
		convSumY := 0.
		for j := y0; j <= y1; j++ {
			for k := x0; k <= x1; k++ {
				h1, h2 := sobel(x-k, y-j)
				convSumX += signal(k, j) * h1
				convSumY += signal(k, j) * h2
			}
		}
		return convSumX, convSumY
	}

	mag := image.NewGray(gray.Bounds())

	wg := sync.WaitGroup{}
	cut := 0.1 * 255
	for y := 0; y < mb; y++ {
		//Proccess lines concurrently
		wg.Add(1)
		go func(y int, mag *image.Gray, wg *sync.WaitGroup) {
			for x := 0; x < nb; x++ {
				convX, convY := conv(x, y)
				value := rescale(math.Hypot(convX, convY), 0, 1530, 0, 255)
				if value <= cut {
					continue
				}
				mag.Pix[y*nb+x] = uint8(value)
			}
			wg.Done()
		}(y, mag, &wg)
	}
	wg.Wait()
	return mag
}
