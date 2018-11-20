package vision

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"

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

// HoughSpace stores the relevant data from the Hough transform
// algorithm.
type HoughSpace struct {
	ThetaAxisSize  int
	RAxisSize      int
	Points         map[int](map[int]int)
	TotalPoints    int
	MaxScore       int
	MinScore       int
	OriginalBounds image.Rectangle
	Data           map[int]*HoughPoint
}

func (h HoughPoint) String() string {
	if len(h.Indexes) < 2 {
		return "{}"
	}
	return fmt.Sprintf("{%4d,%4d:%4d}", h.Indexes[0], h.Indexes[1], h.Score)
}

// Encode returns a unique integer representation for each
// Hough point.
func (h HoughPoint) Encode() int {
	return cantorPairing(h.Indexes[0], h.Indexes[1])
}

type Line struct {
	Slope float64
	Point image.Point
}

func (l Line) String() string {
	return fmt.Sprintf("(Slope: %.2f, X: %4d, Y: %4d)", l.Slope, l.Point.X, l.Point.Y)
}

/* func Hough(input *image.Gray, maxLines int) (houghSpace *image.Gray, lines []*Line) {
	//Sizes of the axis in theta-r space
	const thetaAxisSize = 460
	const rAxisSize = 360 //Must be even

	var houghStructure [rAxisSize][thetaAxisSize]int

	b := input.Bounds()

	//Increments in theta-r
	rmax := math.Hypot(float64(b.Dx()), float64(b.Dy()))
	dr := rmax / (rAxisSize / 2)
	dtheta := math.Pi / float64(thetaAxisSize)
	houghSpace = image.NewGray(image.Rect(0, 0, thetaAxisSize, rAxisSize))
	maxScore := 0
	for x := 0; x < b.Dx(); x++ {
		for y := 0; y < b.Dy(); y++ {
			c := input.GrayAt(x, y)
			if c.Y != 255 {
				continue
			}
			for thetaIndex := 0; thetaIndex < thetaAxisSize; thetaIndex++ {
				theta := dtheta * float64(thetaIndex)
				r := float64(x)*math.Cos(theta) + float64(y)*math.Sin(theta)
				rIndex := rAxisSize/2 - int(math.Floor(r/dr+0.5))
				houghStructure[rIndex][thetaIndex] = houghStructure[rIndex][thetaIndex] + 1
				if houghStructure[rIndex][thetaIndex] > maxScore {
					maxScore = houghStructure[rIndex][thetaIndex]
				}
			}
		}
	}
	lines = make([]*Line, 0)
	return houghSpace, lines
} */

func NewHoughSpace(input *image.Gray, thetaAxisSize, rAxisSize int) *HoughSpace {
	points := map[int]map[int]int{}
	data := map[int]*HoughPoint{}
	b := input.Bounds()
	rmax := math.Hypot(float64(b.Dx()), float64(b.Dy()))
	dr := rmax / (float64(rAxisSize) / 2)
	dtheta := math.Pi / float64(thetaAxisSize)
	maxScore := 0
	minScore := math.MaxInt32
	totalPoints := 0
	for x := 0; x < b.Dx(); x++ {
		for y := 0; y < b.Dy(); y++ {
			c := input.GrayAt(x, y)
			if c.Y != 255 {
				continue
			}
			for thetaIndex := 0; thetaIndex < thetaAxisSize; thetaIndex++ {
				theta := dtheta * float64(thetaIndex)
				r := float64(x)*math.Cos(theta) + float64(y)*math.Sin(theta)
				rIndex := rAxisSize/2 - int(math.Floor(r/dr+0.5))
				if len(points[thetaIndex]) == 0 {
					points[thetaIndex] = map[int]int{}
				}
				score := points[thetaIndex][rIndex] + 1
				points[thetaIndex][rIndex] = score
				totalPoints++
				if score > maxScore {
					maxScore = score
				}
				if score < minScore {
					minScore = score
				}

				//Under test
				key := cantorPairing(thetaIndex, rIndex)
				if _, ok := data[key]; !ok {
					data[key] = &HoughPoint{
						Score:      score,
						SpatialMin: []int{x, y},
						SpatialMax: []int{x, y},
						Indexes:    []int{thetaIndex, rIndex},
					}
				} else {
					p := data[key]
					p.Score = score
					if p.SpatialMin[0] < x || p.SpatialMin[1] < y {
						p.SpatialMin = []int{x, y}
					} else if p.SpatialMax[0] > x || p.SpatialMax[1] > y {
						p.SpatialMax = []int{x, y}
					}
					data[key] = p
				}
			}
		}
	}
	return &HoughSpace{
		ThetaAxisSize:  thetaAxisSize,
		RAxisSize:      rAxisSize,
		Points:         points,
		MaxScore:       maxScore,
		MinScore:       minScore,
		OriginalBounds: b,
		TotalPoints:    totalPoints,
		Data:           data,
	}
}

func (h *HoughSpace) HoughImage() *image.Gray {
	b := image.Rect(0, 0, h.ThetaAxisSize, h.RAxisSize)
	i := image.NewGray(b)
	for thetaIndex := range h.Points {
		for rIndex := range h.Points[thetaIndex] {
			col := color.Gray{Y: uint8(rescale(float64(h.Points[thetaIndex][rIndex]), 0, float64(h.MaxScore), 0, 255))}
			i.SetGray(thetaIndex, rIndex, col)
		}
	}
	return i
}

// LowPassScores filters out all the Hough points that are
// not lesser than the given ratio of the maximum score.
func (h *HoughSpace) LowPassScores(ratio float64) *HoughSpace {
	h2 := &HoughSpace{
		OriginalBounds: h.OriginalBounds,
		RAxisSize:      h.RAxisSize,
		ThetaAxisSize:  h.ThetaAxisSize,
		MinScore:       h.MinScore,
	}
	list := make(TripleList, h.TotalPoints)
	i := 0
	for thetaIndex := range h.Points {
		for rIndex := range h.Points[thetaIndex] {
			score := h.Points[thetaIndex][rIndex]
			list[i] = Triple{thetaIndex, rIndex, score}
			i++
		}
	}
	sort.Sort(list)
	points := map[int](map[int]int){}
	totalPoints := 0
	cutScore := int(ratio * float64(h.MaxScore))
	maxScore := 0
	for _, triple := range list {
		score := triple.Score
		if score >= cutScore {
			break
		}
		totalPoints++
		if len(points[triple.ThetaIndex]) == 0 {
			points[triple.ThetaIndex] = map[int]int{}
		}
		points[triple.ThetaIndex][triple.RIndex] = score
		if score > maxScore {
			maxScore = score
		}
	}
	h2.Points = points
	h2.TotalPoints = totalPoints
	h2.MaxScore = maxScore
	return h2
}

// HighPassScores filters out all the Hough points that are
// not greater than the given ratio of the maximum score.
func (h *HoughSpace) HighPassScores(ratio float64) *HoughSpace {
	h2 := &HoughSpace{
		OriginalBounds: h.OriginalBounds,
		RAxisSize:      h.RAxisSize,
		ThetaAxisSize:  h.ThetaAxisSize,
		MaxScore:       h.MaxScore,
	}
	list := make(TripleList, h.TotalPoints)
	i := 0
	for thetaIndex := range h.Points {
		for rIndex := range h.Points[thetaIndex] {
			score := h.Points[thetaIndex][rIndex]
			list[i] = Triple{thetaIndex, rIndex, score}
			i++
		}
	}
	sort.Sort(list)
	points := map[int](map[int]int){}
	totalPoints := 0
	cutScore := int(ratio * float64(h.MaxScore))
	minScore := math.MaxInt32
	for i := len(list) - 1; i >= 0; i-- {
		triple := list[i]
		score := triple.Score
		if score < cutScore {
			break
		}
		totalPoints++
		if len(points[triple.ThetaIndex]) == 0 {
			points[triple.ThetaIndex] = map[int]int{}
		}
		points[triple.ThetaIndex][triple.RIndex] = score
		if score < minScore {
			minScore = score
		}
	}
	h2.Points = points
	h2.TotalPoints = totalPoints
	h2.MinScore = minScore
	return h2
}

// BandPassScores filters out all the Hough points that are
// greater or lesser than the given ratios of the maximum
// score.
func (h *HoughSpace) BandPassScores(lowerRatio, upperRatio float64) *HoughSpace {
	h2 := &HoughSpace{
		OriginalBounds: h.OriginalBounds,
		RAxisSize:      h.RAxisSize,
		ThetaAxisSize:  h.ThetaAxisSize,
	}
	list := make(TripleList, h.TotalPoints)
	i := 0
	for thetaIndex := range h.Points {
		for rIndex := range h.Points[thetaIndex] {
			score := h.Points[thetaIndex][rIndex]
			list[i] = Triple{thetaIndex, rIndex, score}
			i++
		}
	}
	sort.Sort(list)
	points := map[int](map[int]int){}
	totalPoints := 0
	lowerCutScore := int(lowerRatio * float64(h.MaxScore))
	upperCutScore := int(upperRatio * float64(h.MaxScore))
	minScore := math.MaxInt32
	maxScore := 0
	for i := 0; i < len(list); i++ {
		triple := list[i]
		score := triple.Score
		if score < lowerCutScore {
			continue
		}
		if score > upperCutScore {
			break
		}
		totalPoints++
		if len(points[triple.ThetaIndex]) == 0 {
			points[triple.ThetaIndex] = map[int]int{}
		}
		points[triple.ThetaIndex][triple.RIndex] = score
		if score < minScore {
			minScore = score
		}
		if score > maxScore {
			maxScore = score
		}
	}
	h2.Points = points
	h2.TotalPoints = totalPoints
	h2.MinScore = minScore
	h2.MaxScore = maxScore
	return h2
}

// FindCentroids produces blobs by thresholding an image
// representation of the Hough space and then represents each
// blob with the point closest to its centroid.
func (h *HoughSpace) FindCentroids(threshold uint8) *HoughSpace {
	img := h.HoughImage()
	aux := image.Image(img)
	aux = Threshold(&aux, threshold)
	blobs := *ListBlobs(&aux, Connectivity8)
	h2 := HoughSpace{
		OriginalBounds: h.OriginalBounds,
		ThetaAxisSize:  h.ThetaAxisSize,
		RAxisSize:      h.RAxisSize,
	}
	points := map[int](map[int]int){}
	data := map[int]*HoughPoint{}
	totalPoints := 0
	maxScore := 0
	minScore := math.MaxInt32
	for _, blob := range blobs {
		triple := NewTriple(0, 0, 0)
		totalPoints++
		for _, pt := range blob.Points {
			triple.Score += h.Points[pt.X][pt.Y]
			triple.ThetaIndex += pt.X
			triple.RIndex += pt.Y
		}
		triple.ThetaIndex = int(float64(triple.ThetaIndex) / float64(len(blob.Points)))
		triple.RIndex = int(float64(triple.RIndex) / float64(len(blob.Points)))
		if triple.Score < minScore {
			minScore = triple.Score
		} else if triple.Score > maxScore {
			maxScore = triple.Score
		}
		preAllocMap(&points, triple.ThetaIndex)
		points[triple.ThetaIndex][triple.RIndex] = triple.Score
		//Under test
		key := cantorPairing(triple.ThetaIndex, triple.RIndex)
		if _, ok := h.Data[key]; !ok {
			p := blob.ClosestPoint(blob.Centroid)
			key = cantorPairing(p.X, p.Y)
			h.Data[key].Score = triple.Score
		} else {
			h.Data[key].Score = triple.Score
		}
		data[key] = h.Data[key]
	}
	if maxScore == 0 {
		maxScore = minScore
	}
	h2.Points = points
	h2.MaxScore = maxScore
	h2.MinScore = minScore
	h2.TotalPoints = totalPoints
	h2.Data = data
	return &h2
}

type Triple struct {
	ThetaIndex int
	RIndex     int
	Score      int
}

func NewTriple(thetaIndex, rIndex, score int) Triple {
	return Triple{
		ThetaIndex: thetaIndex,
		RIndex:     rIndex,
		Score:      score,
	}
}

type TripleList []Triple

func (l TripleList) Len() int {
	return len(l)
}

func (l TripleList) Less(i, j int) bool {
	return l[i].Score < l[j].Score
}

func (l TripleList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (h *HoughSpace) GetHighest(n int) *HoughSpace {
	h2 := &HoughSpace{
		MaxScore:       h.MaxScore,
		OriginalBounds: h.OriginalBounds,
		RAxisSize:      h.RAxisSize,
		ThetaAxisSize:  h.ThetaAxisSize,
		TotalPoints:    n,
	}
	list := make(TripleList, h.TotalPoints)
	i := 0
	for thetaIndex := range h.Points {
		for rIndex := range h.Points[thetaIndex] {
			score := h.Points[thetaIndex][rIndex]
			list[i] = Triple{thetaIndex, rIndex, score}
			i++
		}
	}
	sort.Sort(list)
	points := map[int](map[int]int){}
	for _, triple := range list[len(list)-n : len(list)] {
		if len(points[triple.ThetaIndex]) == 0 {
			points[triple.ThetaIndex] = map[int]int{}
		}
		points[triple.ThetaIndex][triple.RIndex] = triple.Score

	}
	h2.Points = points
	return h2
}

func (h *HoughSpace) PlotLines() *image.Gray {
	b := h.OriginalBounds
	context := gg.NewContext(b.Dx(), b.Dy())
	context.SetRGB(0, 0, 0)
	context.Clear()
	context.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{255, 255, 255, 255}))
	context.SetLineWidth(0.4)
	/* rmax := math.Hypot(float64(b.Dx()), float64(b.Dy()))
	dr := rmax / (float64(h.RAxisSize) / 2)
	dtheta := math.Pi / float64(h.ThetaAxisSize) */
	/* for thetaIndex := range h.Points {
		for rIndex := range h.Points[thetaIndex] {
			theta := dtheta * float64(thetaIndex)
			r := float64(h.RAxisSize/2-rIndex) * dr
			x0 := r * math.Cos(theta)
			y0 := r * math.Sin(theta)
			m := math.Tan(theta + math.Pi/2)

			y1 := float64(0)
			x1 := (m*x0 - y0) / m

			if x1 < 0 {
				x1 = float64(b.Dx()) - 1
				y1 = m*(x1-x0) + y0
			}

			x2 := float64(0)
			y2 := -m*x0 + y0

			if y2 < 0 {
				y2 = float64(b.Dy()) - 1
				x2 = (y2-y0)/m + x0
			}

			context.DrawLine(x1, y1, x2, y2)
			context.Stroke()
		}
	} */
	for _, p := range h.Data {
		xmin := float64(p.SpatialMin[0])
		ymin := float64(p.SpatialMin[1])
		xmax := float64(p.SpatialMax[0])
		ymax := float64(p.SpatialMax[1])
		context.DrawLine(xmin, ymin, xmax, ymax)
		context.Stroke()
	}
	/* for thetaIndex := range h.Points {
		rmin := math.MaxFloat64
		rmax := float64(0)
		var tmin, tmax float64
		for rIndex := range h.Points[thetaIndex] {
			if h.Points[thetaIndex][rIndex] == 1 {
				continue
			}
			theta := dtheta * float64(thetaIndex)
			r := float64(h.RAxisSize/2-rIndex) * dr
			if r < rmin {
				rmin = r
				tmin = theta
			} else if r > rmax {
				rmax = r
				tmax = theta
			}
		}
		xmin := rmin * math.Cos(tmin)
		ymin := rmin * math.Sin(tmin)
		xmax := rmax * math.Cos(tmax)
		ymax := rmax * math.Sin(tmax)
		context.DrawLine(xmin, ymin, xmax, ymax)
		context.Stroke()
	} */
	img := context.Image()
	gray := image.NewGray(b)
	draw.Draw(gray, b, img, image.ZP, draw.Src)
	return gray
}

func (h *HoughSpace) ExtractLines() *[]Line {
	return nil
}
