package vision

import (
	"image"
	"image/draw"
)

func Threshold(img *image.Image, level uint8) (out *image.Gray) {
	b := (*img).Bounds()
	out = image.NewGray(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, b, *img, b.Min, draw.Src)
	for i, v := range out.Pix {
		if v <= level {
			out.Pix[i] = 0
		} else {
			out.Pix[i] = 255
		}
	}
	return
}
