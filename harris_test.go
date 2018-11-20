package vision

import (
	"image"
	"image/draw"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func Test_siiGConv(t *testing.T) {
	img, _ := imgio.Open("input_0_sel.png")
	b := img.Bounds()
	gray := image.NewGray(b)
	draw.Draw(gray, b, img, b.Min, draw.Src)
	input := Gray2Mat(gray)
	output := siiGConv(input, 1, 3)
	gray = Mat2Gray(output)
	_ = imgio.Save("SII.png", gray, imgio.PNGEncoder())
}
