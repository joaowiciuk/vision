package vision

import (
	"fmt"
	"image"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func TestBlobs(t *testing.T) {
	img, _ := imgio.Open("examples/Pentagon.png")
	thresholded := image.Image(Threshold(&img, uint8(75)))
	output := Blobs(&thresholded, Connectivity8)
	_ = imgio.Save("examples/Pentagon-blobs.png", output, imgio.PNGEncoder())
}

func TestBlob_ClosestPoint(t *testing.T) {
	img, _ := imgio.Open("examples/blob.png")
	blobs := *ListBlobs(&img, Connectivity8)
	fmt.Println(blobs[0].ClosestPoint(image.Pt(100, 200)))

}
