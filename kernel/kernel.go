//Package kernel provides helper functions for generating convolution kernels.
package kernel

import (
	"math"

	"github.com/joaowiciuk/matrix"
)

// Laplacian generates the laplatian kernel, commonly used for edge detection.
func Laplacian() *matrix.Matrix {
	return &matrix.Matrix{
		{1, 1, 1},
		{1, -8, 1},
		{1, 1, 1},
	}
}

// Sharpen generates the sharpen kernel, used for image enhancement.
func Sharpen() *matrix.Matrix {
	return &matrix.Matrix{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
}

// Line180 generates the 180 degrees edge detector.
func Line180() *matrix.Matrix {
	return &matrix.Matrix{
		{-1, -1, -1},
		{2, 2, 2},
		{-1, -1, -1},
	}
}

// Line90 generates the 90 degrees edge detector.
func Line90() *matrix.Matrix {
	return &matrix.Matrix{
		{-1, 2, -1},
		{-1, 2, -1},
		{-1, 2, -1},
	}
}

// Line45 generates the 45 degrees edge detector.
func Line45() *matrix.Matrix {
	return &matrix.Matrix{
		{2, -1, -1},
		{-1, 2, -1},
		{-1, -1, 2},
	}
}

// Line135 generates the 35 degrees edge detector.
func Line135() *matrix.Matrix {
	return &matrix.Matrix{
		{-1, -1, 2},
		{-1, 2, -1},
		{2, -1, -1},
	}
}

// LoG generates the laplatian of gaussian kernel, commonly used for edge detection.
func LoG() *matrix.Matrix {
	return &matrix.Matrix{
		{0, 0, -1, 0, 0},
		{0, -1, -2, -1, 0},
		{-1, -2, 16, -2, -1},
		{0, -1, -2, -1, 0},
		{0, 0, -1, 0, 0},
	}
}

// Box generates the box blur kernel, used for image blurring.
func Box() *matrix.Matrix {
	A := &matrix.Matrix{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	}
	return A.Scal(1. / 9.)
}

// Unsharp55 generates the 5x5 unsharp kernel.
func Unsharp55() *matrix.Matrix {
	A := &matrix.Matrix{
		{1, 4, 6, 4, 1},
		{4, 16, 24, 16, 4},
		{6, 24, -476, 24, 6},
		{4, 16, 24, 16, 4},
		{1, 4, 6, 4, 1},
	}
	return A.Scal(-1. / 256.)
}

// SobelX generates the n-by-n horizontal sobel edge detector.
func SobelX(n int) *matrix.Matrix {
	if n < 3 {
		n = 3
	}
	if n%2 == 0 {
		n++
	}
	A := matrix.New(n, n)
	A.Law(func(r, c int) float64 {
		y := float64(c-n/2) / (math.Pow(float64(c-n/2), 2) + math.Pow(float64(r-n/2), 2))
		if !math.IsNaN(y) {
			return y
		}
		return 0
	})
	return A
}

// SobelY generates the n-by-n vertical sobel edge detector.
func SobelY(n int) *matrix.Matrix {
	if n < 3 {
		n = 3
	}
	if n%2 == 0 {
		n++
	}
	A := matrix.New(n, n)
	A.Law(func(r, c int) float64 {
		y := float64(r-n/2) / (math.Pow(float64(c-n/2), 2) + math.Pow(float64(r-n/2), 2))
		if !math.IsNaN(y) {
			return y
		}
		return 0
	})
	return A
}

// Gaussian generates the n-by-n gaussian kernel with standart deviation σ, commonly used for image blurring.
func Gaussian(n int, σ float64) *matrix.Matrix {
	if n < 3 {
		n = 3
	}
	if n%2 == 0 {
		n++
	}
	X := matrix.New(n, n)
	u := float64(n / 2)
	X.Law(func(r, c int) float64 {
		x := float64(c)
		y := float64(r)
		return math.Exp(-(math.Pow(x-u, 2) + math.Pow(y-u, 2)) / (2. * σ * σ))
	})
	s := X.Sum()
	return X.ForEach(func(x float64) float64 { return x / s })
}

// Gaussian1D generates the 1D gaussian kernel with standart deviation σ, used for separable convolution.
func Gaussian1D(n int, σ float64) (A *matrix.Matrix) {
	if n < 3 {
		n = 3
	}
	if n%2 == 0 {
		n++
	}
	X := matrix.New(n, n)
	u := float64(n / 2)
	X.Law(func(r, c int) float64 {
		x := float64(c)
		y := float64(r)
		return math.Exp(-(math.Pow(x-u, 2) + math.Pow(y-u, 2)) / (2. * σ * σ))
	})
	s := X.Sum()
	X = X.ForEach(func(x float64) float64 { return x / s })
	return X.Row(n / 2)
}
