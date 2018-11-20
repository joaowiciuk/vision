package vision

import (
	"image"
	"image/draw"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func TestGrad(t *testing.T) {
	img, _ := imgio.Open("examples/rune.png")
	b := img.Bounds()
	gray := image.NewGray(b)
	draw.Draw(gray, b, img, b.Min, draw.Src)
	mag, ang := Grad(gray)
	_ = imgio.Save("mag.png", mag, imgio.PNGEncoder())
	_ = imgio.Save("ang.png", ang, imgio.PNGEncoder())
}
