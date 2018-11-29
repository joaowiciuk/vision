package vision

import (
	"image"
	"image/draw"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func TestGaussian(t *testing.T) {
	img, _ := imgio.Open("images/house.png")
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.ZP, draw.Src)
	gaussian := Gaussian(gray, 1.4)
	_ = imgio.Save("images/house-gaussian.png", gaussian, imgio.PNGEncoder())
}
