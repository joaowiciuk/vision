package vision

import (
	"fmt"
	"image"
	"testing"

	"github.com/joaowiciuk/lenna/io"
)

func TestStack(t *testing.T) {
	img := *io.LoadRGBA("/home/joaowiciuk/Imagens/lenna.jpg")
	R := image.Rect(0, 0, 128, 512)
	O := Stack(img, R)
	/* if len(O) != 9 {
		t.Fatal("Wrong number of frames")
	} */
	for i := range O {
		if O[i].Bounds().Dx() != R.Bounds().Dx() || O[i].Bounds().Dy() != R.Bounds().Dy() {
			t.Fatal("Wrong frame")
		}
		io.SavePNG(O[i], fmt.Sprintf("/home/joaowiciuk/Imagens/lenna.jpg (%d).png", i))
	}
}
