package vision

import (
	"fmt"
	"image"
	"math"

	"github.com/joaowiciuk/matrix"
)

// Harris implements the Harris corner detector as described in
//Javier Sánchez, Nelson Monzón, and Agustín Salgado,
//An Analysis and Implementation of the Harris Corner Detector,
//Image Processing On Line, 8 (2018), pp. 305–328.
//https://doi.org/10.5201/ipol.2018.229
func Harris(img *image.Gray, measure, k, d, i float32, strategy string, cells, N int, subpixel bool) (x, y *[]int) {
	//Smoothing the image
	/* mat := Gray2Mat(img) */
	//Computing the gradient of the image
	//Computing the autocorrelation matrix
	//Non-maximum suppression
	//Selecting output corners
	//Calculating subpixel accuracy
	return nil, nil
}

func siiGConv(input *matrix.Matrix, d float64, k int) *matrix.Matrix {
	//Select the initial radii and weights based on k
	/* fmt.Println("Selecting initial raddi and weights...")
	var r0 []float64
	var w0 []float64
	switch k {
	case 3:
		r0 = []float64{23, 46, 76}
		w0 = []float64{0.9495, 0.5502, 0.1618}
	case 4:
		r0 = []float64{19, 37, 56, 82}
		w0 = []float64{0.9649, 0.6700, 0.3376, 0.0976}
	case 5:
		r0 = []float64{16, 30, 44, 61, 85}
		w0 = []float64{0.9738, 0.7596, 0.5031, 0.2534, 0.0739}
	} */

	//Adjust the radii and weights for the given standart deviation d
	/* fmt.Println("Adjusting raddi and weights using the standart deviation...")
	r := make([]float64, k)
	w := make([]float64, k)
	d0 := float64(100 / math.Pi)
	rmax := float64(0)
	for i := 0; i < k; i++ {
		r[i] = d / d0 * r0[i]
		if r[i] > rmax {
			rmax = r[i]
		}
	}
	fmt.Printf("radii adjusted and max radius = %.4f\n", rmax)
	for i := 0; i < k; i++ {
		sum := float64(0)
		for j := 0; j < k; j++ {
			sum += w0[j] * (2*r[j] + 1)
		}
		w[i] = w0[i] / sum
	}

	fmt.Println("Actual radii and weights:")
	fmt.Printf("Radii %v\n", r)
	fmt.Printf("Weights %v\n", w) */

	/* fmt.Println("Newest radii and weights:") */
	radii, weights := siiCoeffs(d, k)
	fmt.Printf("Radii %v\n", radii)
	fmt.Printf("Weights %v\n", weights)

	rows, cols := input.Size()

	//Compute the horizontal cumulative sums
	fmt.Println("Computing the horizontal cumulative sums...")
	hsums := matrix.New(rows, cols)
	for y := 0; y < rows; y++ {
		for x := 1; x < cols; x++ {
			hsums.Set(y, x, input.At(y, x)+input.At(y, x-1))
		}
	}

	output := matrix.New(rows, cols)

	//Compute the horizontal convolution
	fmt.Println("Computing the horizontal convolution...")
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			sum := float64(0)
			for i := 0; i < k; i++ {
				var index1, index2 int
				if x+int(radii[i]) < cols {
					index1 = x + int(radii[i])
				} else {
					index1 = cols - 1
				}
				if x-int(radii[i])-1 >= 0 {
					index2 = x - int(radii[i]) - 1
				} else {
					index2 = 0
				}
				sum += weights[i] * (hsums.At(y, index1) - hsums.At(y, index2))
			}
			output.Set(y, x, sum)
		}
	}

	//Compute the vertical cumulative sums (now based on previous convolution)
	fmt.Println("Computing the vertical cumulative sums...")
	vsums := matrix.New(rows, cols)
	for x := 0; x < cols; x++ {
		for y := 1; y < rows; y++ {
			vsums.Set(y, x, output.At(y, x)+output.At(y-1, x))
		}
	}

	//Compute the vertical convolution (now based on previous convolution)
	fmt.Println("Computing the vertical convolution...")
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			sum := float64(0)
			for i := 0; i < k; i++ {
				var index1, index2 int
				if y+int(radii[i]) < rows {
					index1 = y + int(radii[i])
				} else {
					index1 = rows - 1
				}
				if y-int(radii[i])-1 >= 0 {
					index2 = y - int(radii[i]) - 1
				} else {
					index2 = 0
				}
				sum += weights[i] * (vsums.At(index1, x) - vsums.At(index2, x))
			}
			output.Set(y, x, sum)
		}
	}

	return output
}

func siiCoeffs(sigma float64, k int) ([]int, []float64) {
	const minK = 3
	const maxK = 5
	sigma0 := 100 / math.Pi
	radii0 := [maxK - minK + 1][maxK]int{
		{76, 46, 23, 0, 0},
		{82, 56, 37, 19, 0},
		{85, 61, 44, 30, 16},
	}
	weights0 := [maxK - minK + 1][maxK]float64{
		{0.1618, 0.5502, 0.9495, 0, 0},
		{0.0976, 0.3376, 0.6700, 0.9649, 0},
		{0.0739, 0.2534, 0.5031, 0.7596, 0.9738},
	}
	i := k - minK
	sum := float64(0)
	radii := make([]int, k)
	weights := make([]float64, k)
	for j := 0; j < k; j++ {
		radii[j] = int(float64(radii0[i][j]) * sigma / sigma0)
		sum += weights0[i][j] * (2*float64(radii[j]) + 1)
	}
	for j := 0; j < k; j++ {
		weights[j] = weights0[i][j] / sum
	}
	return radii, weights
}
