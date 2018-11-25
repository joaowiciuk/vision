package vision

import (
	"image"
	"math"
)

func Gaussian(gray *image.Gray, sigma float64) *image.Gray {
	var c sii_coeffs
	sii_precomp(&c, sigma, 3)
	buffer := make([]float64, sii_buffer_size(c, len(gray.Pix)))
	output_image := make([]float64, len(gray.Pix))
	input_image := make([]float64, len(gray.Pix))
	for i, p := range gray.Pix {
		input_image[i] = float64(p)
	}
	width := gray.Bounds().Dx()
	height := gray.Bounds().Dy()
	var output_image_i, buffer_i, input_image_i int
	sii_gaussian_conv_image(c, &output_image, &output_image_i, &buffer, &buffer_i, &input_image, &input_image_i, width, height, 1)
	outputGray := image.NewGray(gray.Bounds())
	for i := range output_image {
		outputGray.Pix[i] = uint8(output_image[i])
	}
	return outputGray
}

type sii_coeffs struct {
	weights []float64
	radii   []int
	K       int
}

func sii_precomp(c *sii_coeffs, sigma float64, K int) {
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
	i := K - minK
	sum := float64(0)
	radii := make([]int, K)
	weights := make([]float64, K)
	for j := 0; j < K; j++ {
		radii[j] = int(float64(radii0[i][j]) * sigma / sigma0)
		sum += weights0[i][j] * (2*float64(radii[j]) + 1)
	}
	for j := 0; j < K; j++ {
		weights[j] = weights0[i][j] / sum
	}
	c.radii = radii
	c.weights = weights
	c.K = K
}

func sii_buffer_size(c sii_coeffs, N int) int {
	pad := c.radii[0] + 1
	return N + 2*pad
}

func sii_gaussian_conv(c sii_coeffs, dest *[]float64, dest_i *int, buffer *[]float64, buffer_i *int, src *[]float64, N int, stride int) {
	var accum float64
	var pad int
	pad = c.radii[0] + 1
	(*buffer_i) += pad

	for n := -pad; n < N+pad; n++ {
		accum += (*src)[stride*extension(N, n)]
		(*buffer)[n] = accum
	}
	for n := 0; n < N; n++ {
		accum = c.weights[0] * float64((*buffer)[n+c.radii[0]]-(*buffer)[n-c.radii[0]-1])
		for k := 1; k < c.K; k++ {
			accum += c.weights[k] * float64((*buffer)[n+c.radii[k]]-(*buffer)[n-c.radii[k]-1])
		}
		(*dest_i) += stride
	}
}

func extension(N, n int) int {
	for {
		if n < 0 {
			n = -1 - n
		} else if n >= N {
			n = 2*N - 1 - 1
		} else {
			break
		}
	}
	return n
}

func sii_gaussian_conv_image(c sii_coeffs, dest *[]float64, dest_i *int, buffer *[]float64, buffer_i *int, src *[]float64, src_i *int, width, height, num_channels int) {
	num_pixels := width * height
	/* Loop over the image channels. */
	for channel := 0; channel < num_channels; channel++ {
		dest_y := dest
		src_y := src

		var dest_y_i, src_y_i *int
		/* (*dest_y_i) = 0
		(*src_y_i) = 0 */

		/* Filter each row of the channel. */
		for y := 0; y < height; y++ {
			sii_gaussian_conv(c, dest_y, dest_y_i, buffer, buffer_i, src_y, width, 1)
			(*dest_y_i) += width
			(*src_y_i) += width
		}

		/* Filter each column of the channel. */
		for x := 0; x < width; x++ {
			(*dest_i) += x
			sii_gaussian_conv(c, dest, dest_i, buffer, buffer_i, dest, height, width)
		}
		(*dest_i) += num_pixels
		(*src_i) += num_pixels
	}

	return
}
