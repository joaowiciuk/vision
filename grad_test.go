package vision

import (
	"image"
	"image/draw"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func TestGrad(t *testing.T) {
	img, _ := imgio.Open("images/house.png")
	b := img.Bounds()
	gray := image.NewGray(b)
	draw.Draw(gray, b, img, b.Min, draw.Src)
	mag, ang := Grad(gray)
	_ = imgio.Save("images/house-mag.png", mag, imgio.PNGEncoder())
	_ = imgio.Save("images/house-ang.png", ang, imgio.PNGEncoder())
}

func BenchmarkGrad(b *testing.B) {
	b.StopTimer()
	img, _ := imgio.Open("images/house.png")
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, img, bounds.Min, draw.Src)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Grad(gray)
	}
}
