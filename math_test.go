package vision

import (
	"fmt"
	"image"
	"image/draw"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func Test_generateVoronoi(t *testing.T) {
	img, _ := imgio.Open("examples/blob.png")
	blobs := *ListBlobs(&img, Connectivity8)
	var sx, sy []int
	for _, b := range blobs {
		for _, p := range b.Points {
			sx = append(sx, p.X)
			sy = append(sy, p.Y)
		}
	}
	fmt.Println(sx)
	fmt.Println(sy)
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.ZP, draw.Src)
	voroni := generateVoronoi(gray, sx, sy)
	_ = imgio.Save("examples/blob-voroni.png", voroni, imgio.PNGEncoder())
}
