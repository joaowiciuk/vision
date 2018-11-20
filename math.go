package vision

import (
	"image"
	"image/color"
	"math/rand"
)

func max(i, j int) int {
	if i >= j {
		return i
	}
	return j
}

func min(i, j int) int {
	if i <= j {
		return i
	}
	return j
}

func generateVoronoi(binary *image.Gray, sx, sy []int) image.Image {
	// generate a random color for each site
	nSites := len(sx)
	imageWidth, imageHeight := binary.Bounds().Dx(), binary.Bounds().Dy()
	sc := make([]color.NRGBA, nSites)
	for i := range sx {
		sc[i] = color.NRGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)),
			uint8(rand.Intn(256)), 255}
	}
	// generate diagram by coloring each pixel with color of nearest site
	img := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			dMin := dot(imageWidth, imageHeight)
			var sMin int
			for s := 0; s < nSites; s++ {
				if d := dot(sx[s]-x, sy[s]-y); d < dMin {
					sMin = s
					dMin = d
				}
			}
			img.SetNRGBA(x, y, sc[sMin])
		}
	}
	return img
}

func dot(x, y int) int {
	return x*x + y*y
}
