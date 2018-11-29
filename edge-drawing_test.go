package vision

import (
	"image"
	"image/draw"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func TestEdgeDrawing(t *testing.T) {
	img, _ := imgio.Open("images/house.png")
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.ZP, draw.Src)
	edges := EdgeDrawing(gray)
	_ = imgio.Save("images/house-edges.png", edges, imgio.PNGEncoder())
}

func BenchmarkEdgeDrawing(b *testing.B) {
	b.StopTimer()
	img, _ := imgio.Open("images/house.png")
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.ZP, draw.Src)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = EdgeDrawing(gray)
	}
}
