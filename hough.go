package vision

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/fogleman/gg"
)

// HoughPoint represents a single point in the Hough space with
// its score, theta and rho values and and minimum and maximum
// corresponding spatial points.
type HoughPoint struct {
	Indexes    []int
	Score      int
	SpatialMin []int
	SpatialMax []int
}

// String returns a string representation of the HoughPoint.
func (h HoughPoint) String() string {
	if len(h.Indexes) < 2 {
		return "{}"
	}
	return fmt.Sprintf("{%4d,%4d:%4d}", h.Indexes[0], h.Indexes[1], h.Score)
}

// GenerateKey returns an unique integer representation for each
// Hough point.
func GenerateKey(theta, rho int) int {
	return pairing(theta, rho)
}

// HoughSpace stores the relevant data from the Hough transform
// algorithm.
type HoughSpace struct {
	ThetaRes      int
	RhoRes        int
	MaxScore      int
	MinScore      int
	SpatialBounds image.Rectangle
	points        map[int]*HoughPoint
}

// NewHoughSpace performs the Hough transform in the input space image
// and returns a HoughSpace struct pointer with the given theta and rho
// resolutions.
func NewHoughSpace(input *image.Gray, thetaRes, rhoRes int) *HoughSpace {
	hs := &HoughSpace{
		ThetaRes: thetaRes,
		RhoRes:   rhoRes,
		points:   map[int]*HoughPoint{},
	}
	b := input.Bounds()
	hs.SpatialBounds = b
	rhoMax := math.Hypot(float64(b.Dx()), float64(b.Dy()))
	drho := rhoMax / float64(rhoRes/2)
	/* dtheta := math.Pi / float64(thetaRes) */
	hs.MaxScore = 0
	hs.MinScore = math.MaxInt32
	/* img := image.NewGray(image.Rect(0, 0, thetaRes, rhoRes)) */
	for x := 0; x < b.Dx(); x++ {
		for y := 0; y < b.Dy(); y++ {
			c := input.GrayAt(x, y)
			if c.Y != 255 {
				continue
			}
			for thetaIndex := 0; thetaIndex < thetaRes; thetaIndex++ {
				//Rescale thetaIndex to the interval [0, pi] and assign to theta
				/* var theta float64
				if thetaIndex > 0 && thetaIndex < thetaRes-1 {
					theta = dtheta*(float64(thetaIndex)+math.Pi) + (float64(thetaRes) - 1)
				} else if thetaIndex == 0 {
					theta = 0
				} else if thetaIndex == thetaRes-1 {
					theta = math.Pi
				} */
				theta := rescale(float64(thetaIndex), 0, float64(thetaRes-1), -math.Pi/2, math.Pi/2)
				rho := float64(x)*math.Cos(theta) + float64(y)*math.Sin(theta)
				rhoIndex := rhoRes/2 - int(math.Floor(rho/drho+0.5))
				var score int
				key := pairing(thetaIndex, rhoIndex)
				if _, ok := hs.points[key]; !ok {
					score = 1
					hs.points[key] = &HoughPoint{
						Score:      score,
						SpatialMin: []int{x, y},
						SpatialMax: []int{x, y},
						Indexes:    []int{thetaIndex, rhoIndex},
					}
				} else {
					p := hs.points[key]
					p.Score++
					score = p.Score
					if p.SpatialMin[0] < x || p.SpatialMin[1] < y {
						p.SpatialMin = []int{x, y}
					} else if p.SpatialMax[0] > x || p.SpatialMax[1] > y {
						p.SpatialMax = []int{x, y}
					}
					hs.points[key] = p
				}
				if score < hs.MinScore {
					hs.MinScore = score
				} else if score > hs.MaxScore {
					hs.MaxScore = score
				}

				/* col := img.At(thetaIndex, rhoIndex).(color.Gray)
				if col.Y < 255 {
					col.Y++
					img.SetGray(thetaIndex, rhoIndex, col)
				} */
			}
		}
	}
	/* _ = imgio.Save("hue.png", img, imgio.PNGEncoder()) */
	if hs.MaxScore == 0 {
		hs.MaxScore = hs.MinScore
	}
	return hs
}

// At returns the Hough point at theta and rho indexes.
func (h *HoughSpace) At(theta, rho int) (*HoughPoint, bool) {
	key := GenerateKey(theta, rho)
	hp, ok := h.points[key]
	return hp, ok
}

// Set overwrites the Hough point at theta and rho indexes.
func (h *HoughSpace) Set(theta, rho int, hp *HoughPoint) {
	key := GenerateKey(theta, rho)
	h.points[key] = hp
}

// Count returns the total number of Hough points.
func (h *HoughSpace) Count() int {
	return len(h.points)
}

// HoughImage rescales the scores in the Hough space to
// the interval [0, 255] and plot it in a grayscale image.
func (h *HoughSpace) HoughImage() *image.Gray {
	b := image.Rect(0, 0, h.ThetaRes, h.RhoRes)
	i := image.NewGray(b)
	for _, point := range h.points {
		col := color.Gray{Y: uint8(rescale(float64(point.Score), 0, float64(h.MaxScore), 0, 255))}
		i.SetGray(point.Indexes[0], point.Indexes[1], col)
	}
	return i
}

// FindCentroids thresholds the image of the Hough space,
// find its blobs and return a new Hough space containing
// only the most representative point for each blob.
func (h *HoughSpace) FindCentroids(threshold uint8) *HoughSpace {
	img := h.HoughImage()
	aux := image.Image(img)
	aux = Threshold(&aux, threshold)
	blobs := *ListBlobs(&aux, Connectivity8)
	h2 := &HoughSpace{
		SpatialBounds: h.SpatialBounds,
		ThetaRes:      h.ThetaRes,
		RhoRes:        h.RhoRes,
		points:        map[int]*HoughPoint{},
		MaxScore:      0,
		MinScore:      math.MaxInt32,
	}
	for _, blob := range blobs {
		//Compute the blob centroid
		var thetaMomentum, rhoMomentum float64
		var totalScore int64
		for _, p := range blob.Points {
			score := h.points[pairing(p.X, p.Y)].Score
			totalScore += int64(score)
			thetaMomentum += float64(p.X * score)
			rhoMomentum += float64(p.Y * score)
		}
		theta := int(thetaMomentum / float64(totalScore))
		rho := int(rhoMomentum / float64(totalScore))

		//Checks if the centroid lies over a point
		var hp *HoughPoint
		key := pairing(theta, rho)
		if hp0, ok := h.points[key]; ok {
			//If it does, use it for representing the blob
			hp = hp0
		} else {
			//If not, use the nearest point that does
			p := blob.ClosestPoint(image.Pt(theta, rho))
			key = pairing(p.X, p.Y)
			hp = h.points[key]
		}

		//Store the point in the new Hough space
		h2.points[key] = hp

		//Update minimum and maximum scores
		if hp.Score < h2.MinScore {
			h2.MinScore = hp.Score
		} else if hp.Score > h2.MaxScore {
			h2.MaxScore = hp.Score
		}
	}
	if h2.MaxScore == 0 {
		h2.MaxScore = h2.MinScore
	}
	return h2
}

// PlotLines returns an image with line segments for each
// corresponding point in the Hough space.
func (h *HoughSpace) PlotLines() *image.Gray {
	b := h.SpatialBounds
	context := gg.NewContext(b.Dx(), b.Dy())
	context.SetRGB(0, 0, 0)
	context.Clear()
	context.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{255, 255, 255, 255}))
	context.SetLineWidth(0.4)
	for _, p := range h.points {
		xmin := float64(p.SpatialMin[0])
		ymin := float64(p.SpatialMin[1])
		xmax := float64(p.SpatialMax[0])
		ymax := float64(p.SpatialMax[1])
		context.DrawLine(xmin, ymin, xmax, ymax)
		context.Stroke()
	}
	img := context.Image()
	gray := image.NewGray(b)
	draw.Draw(gray, b, img, image.ZP, draw.Src)
	return gray
}
