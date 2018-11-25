package vision

import (
	"image"
)

func EdgeDrawing(gray *image.Gray) *image.Gray {
	mag, _ := Grad(gray)
	aux := image.Image(mag)
	output := Threshold(&aux, 102) //40%
	return output
}
