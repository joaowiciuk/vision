package vision

import (
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

func TestThreshold(t *testing.T) {
	img, _ := imgio.Open("examples/input_0_sel.png")
	output := Threshold(&img, uint8(75))
	_ = imgio.Save("examples/input_0_sel_threshold.png", output, imgio.PNGEncoder())
}
