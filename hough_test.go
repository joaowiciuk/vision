package vision

import (
	"fmt"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func TestHoughSpace_FindCentroids(t *testing.T) {
	img, _ := imgio.Open("examples/rune.png")
	edges := Canny(img, 91, 31, 3, 1.7)
	var houghSpace *HoughSpace
	houghSpace = NewHoughSpace(edges, 800, 420)
	fmt.Println("Max score: ", houghSpace.MaxScore)
	fmt.Println("Min score: ", houghSpace.MinScore)
	fmt.Println("Total points: ", houghSpace.TotalPoints)
	spaceImage := houghSpace.HoughImage()
	houghSpace = houghSpace.FindCentroids(20)
	connectedSpaceImage := houghSpace.HoughImage()
	fmt.Println("Max score after: ", houghSpace.MaxScore)
	fmt.Println("Min score after: ", houghSpace.MinScore)
	fmt.Println("Total points after: ", houghSpace.TotalPoints)
	lines := houghSpace.PlotLines()
	_ = imgio.Save("examples/rune-edges.png", edges, imgio.PNGEncoder())
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
	edges := Canny(img, 91, 31, 3, 0.25)
	var houghSpace *HoughSpace
	houghSpace = NewHoughSpace(edges, 460, 360)
	fmt.Println("Max score: ", houghSpace.MaxScore)
	fmt.Println("Min score: ", houghSpace.MinScore)
	fmt.Println("Total points: ", houghSpace.TotalPoints)
	fmt.Println("Total data: ", len(houghSpace.Data))
}
