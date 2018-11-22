package vision

import (
	"fmt"
	"image"
	"image/draw"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func TestHoughSpace_FindCentroids(t *testing.T) {
	/* img, _ := imgio.Open("examples/shape.png")
	edges := Canny(img, 91, 31, 3, 0.45) */

	img, _ := imgio.Open("examples/rune-edges.png")
	edges := image.NewGray(img.Bounds())
	draw.Draw(edges, img.Bounds(), img, image.ZP, draw.Src)

	var houghSpace *HoughSpace
	houghSpace = NewHoughSpace(edges, 460, 360)
	fmt.Println("Max score: ", houghSpace.MaxScore)
	fmt.Println("Min score: ", houghSpace.MinScore)
	fmt.Println("Total points: ", houghSpace.Count())
	spaceImage := houghSpace.HoughImage()
	houghSpace = houghSpace.FindCentroids(250)
	connectedSpaceImage := houghSpace.HoughImage()
	fmt.Println("Max score after: ", houghSpace.MaxScore)
	fmt.Println("Min score after: ", houghSpace.MinScore)
	fmt.Println("Total points after: ", houghSpace.Count())
	lines := houghSpace.PlotLines()
	/* _ = imgio.Save("examples/rune-edges.png", edges, imgio.PNGEncoder()) */
	_ = imgio.Save("examples/rune-hough-space.png", spaceImage, imgio.PNGEncoder())
	_ = imgio.Save("examples/rune-hough-space-centroids.png", connectedSpaceImage, imgio.PNGEncoder())
	_ = imgio.Save("examples/rune-hough-lines.png", lines, imgio.PNGEncoder())
}

func BenchmarkHoughSpace_PlotLines(b *testing.B) {
	b.StopTimer()
	img, _ := imgio.Open("examples/photo.jpg")
	edges := Canny(img, 91, 31, 3, 1.7)
	var houghSpace *HoughSpace
	houghSpace = NewHoughSpace(edges, 800, 480)
	houghSpace = houghSpace.FindCentroids(200)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = houghSpace.PlotLines()
	}
}

func TestNewHoughSpace(t *testing.T) {
	img, _ := imgio.Open("examples/Pentagon.png")
	/* edges := image.NewGray(img.Bounds())
	draw.Draw(edges, img.Bounds(), img, image.ZP, draw.Src) */
	edges := Canny(img, 91, 31, 3, 0.25)
	var houghSpace *HoughSpace
	houghSpace = NewHoughSpace(edges, 460, 360)
	fmt.Println("Max score: ", houghSpace.MaxScore)
	fmt.Println("Min score: ", houghSpace.MinScore)
	fmt.Println("Total points: ", houghSpace.Count())
	_ = imgio.Save("examples/Pentagon-hough-space.png", houghSpace.HoughImage(), imgio.PNGEncoder())
}

func BenchmarkNewHoughSpace(b *testing.B) {
	b.StopTimer()
	const thetaRes = 460
	const rhoRes = 360
	img, _ := imgio.Open("examples/shape-edges.png")
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, img.Bounds(), img, image.ZP, draw.Src)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = NewHoughSpace(gray, thetaRes, rhoRes)
	}
}
