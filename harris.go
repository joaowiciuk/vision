package vision

import (
	"image"
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
