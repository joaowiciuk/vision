package vision

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Blob represents a blob
type Blob struct {
	Bounds   image.Rectangle
	Centroid image.Point
	Area     int
	Points   []image.Point
}

// Connectivity is an image graph connectivity for use in blob detection
type Connectivity int

const (
	// Connectivity8 8-connected image graph
	Connectivity8 Connectivity = iota

	// Connectivity4 4-connected image graph
	Connectivity4
)

const (
	topLeft  int = 0
	top      int = 1
	topRight int = 2
	left     int = 3
)

// Blobs paint all connected white blobs in a different RGB color
func Blobs(input *image.Image, connectivity Connectivity) *image.RGBA {
	blobs := ListBlobs(input, connectivity)
	imgOut := image.NewRGBA((*input).Bounds())

	//Generate a different random color for each component
	used := make(map[color.RGBA]bool)
	colors := make([]color.RGBA, len(*blobs))
	rand.Seed(time.Now().Unix())
	for i := range *blobs {
		var c color.RGBA
		for {
			c = color.RGBA{
				R: uint8(rand.Intn(255)),
				G: uint8(rand.Intn(255)),
				B: uint8(rand.Intn(255)),
				A: 255,
			}
			if _, ok := used[c]; !ok {
				used[c] = true
				colors[i] = c
				break
			}
		}
	}

	for i, blob := range *blobs {
		for _, point := range blob.Points {
			imgOut.SetRGBA(point.X, point.Y, colors[i])
		}
	}

	//Paint components
	/* if paintComponents {
	} */

	//Paint bounds
	/* if paintBounds {
		for i, blob := range *blobs {
			for col := blob.Bounds.Min.X; col <= blob.Bounds.Max.X; col++ {
				imgOut.SetRGBA(col, blob.Bounds.Min.Y, colors[i])
				imgOut.SetRGBA(col, blob.Bounds.Max.Y, colors[i])
			}
			for row := blob.Bounds.Min.Y; row <= blob.Bounds.Max.Y; row++ {
				imgOut.SetRGBA(blob.Bounds.Min.X, row, colors[i])
				imgOut.SetRGBA(blob.Bounds.Max.X, row, colors[i])
			}
		}
	} */
	//Paint centroids
	/* if paintCentroids {
		for i, blob := range *blobs {
			if !paintComponents {
				imgOut.SetRGBA(blob.Centroid.X, blob.Centroid.Y, colors[i])
			} else {
				inverseColor := color.RGBA{
					R: 255 - colors[i].R,
					G: 255 - colors[i].G,
					B: 255 - colors[i].B,
				}
				imgOut.SetRGBA(blob.Centroid.X, blob.Centroid.Y, inverseColor)
			}
		}
	} */
	return imgOut
}

// ListBlobs returns a list of all white connected blobs
func ListBlobs(i *image.Image, connectivity Connectivity) *[]Blob {
	b := (*i).Bounds()
	img := image.NewGray(b)
	draw.Draw(img, b, *i, b.Min, draw.Src)
	maxBlobs := len(img.Pix)/2 + 2
	blobs := make([]Blob, 0)
	quickUnion := newQuickUnion(maxBlobs)
	width := img.Rect.Size().X
	height := img.Rect.Size().Y
	output := make([][]int, height+2)
	input := make([][]int, height+2)
	count := 0
	neighbor := make([]int, 4)

	one := func(row, col int) bool {
		return input[row][col] == 1
	}

	label := func(row, col int) {
		count++
		output[row][col] = count
	}

	//Zero pass
	input[0] = make([]int, width+2)
	output[0] = make([]int, width+2)
	input[height+1] = make([]int, width+2)
	output[height+1] = make([]int, width+2)
	for row := 1; row < height+1; row++ {
		output[row] = make([]int, width+2)
		input[row] = make([]int, width+2)
		for col := 1; col < width+1; col++ {
			if img.GrayAt(col-1, row-1).Y == 255 {
				input[row][col] = 1
			}
		}
	}

	// First pass
	for row := 1; row < height+1; row++ {
		for col := 1; col < width+1; col++ {

			if one(row, col) {

				switch connectivity {
				case Connectivity4:
					neighbor[top] = output[row-1][col]
					neighbor[left] = output[row][col-1]

					if neighbor[top] == 0 && neighbor[left] == 0 {
						label(row, col)
						continue
					}
					if neighbor[top] == 0 {
						output[row][col] = neighbor[left]
						continue
					}
					if neighbor[left] == 0 {
						output[row][col] = neighbor[top]
						continue
					}
					if neighbor[top] == neighbor[left] {
						output[row][col] = neighbor[top]
						continue
					}
					if neighbor[top] < neighbor[left] {
						output[row][col] = neighbor[top]
						quickUnion.unite(neighbor[top], neighbor[left])
						continue
					}
					if neighbor[top] > neighbor[left] {
						output[row][col] = neighbor[left]
						quickUnion.unite(neighbor[left], neighbor[top])
						continue
					}
				case Connectivity8:
					aux := make([]int, 0)
				outerLoop:
					for row1 := row - 1; row1 <= row+1; row1++ {
						for col1 := col - 1; col1 <= col+1; col1++ {
							if output[row1][col1] != 0 {
								aux = append(aux, output[row1][col1])
							}
							if row1 == row && col1 == col-1 {
								break outerLoop
							}
						}
					}
					if len(aux) == 0 {
						label(row, col)
						continue
					}
					if len(aux) == 1 {
						output[row][col] = aux[0]
						continue
					}
					if len(aux) > 1 {
						sort.Ints(aux)
						minimum := aux[0]
						for _, i := range aux {
							if i != minimum {
								quickUnion.unite(minimum, i)
							}
						}
						output[row][col] = minimum
						continue
					}
				}
			}
		}
	}

	// Second pass
	c := 0
	labelIndex := make(map[int]int)
	for row := 1; row < height+1; row++ {
		for col := 1; col < width+1; col++ {
			if !one(row, col) {
				continue
			}
			label := quickUnion.ID[output[row][col]]
			root := quickUnion.root(label)
			output[row][col] = root
			point := image.Point{X: col - 1, Y: row - 1}
			if pos, ok := labelIndex[root]; !ok {
				labelIndex[root] = c
				blobs = append(blobs, Blob{
					Area:     0,
					Bounds:   image.Rect(point.X-1, point.Y-1, point.X-1, point.Y-1),
					Centroid: point,
					Points:   []image.Point{point},
				})
				c++
			} else {
				blobs[pos].Area++
				if point.X < blobs[pos].Bounds.Min.X {
					blobs[pos].Bounds.Min.X = point.X
				}
				if point.X > blobs[pos].Bounds.Max.X {
					blobs[pos].Bounds.Max.X = point.X
				}
				if point.Y < blobs[pos].Bounds.Min.Y {
					blobs[pos].Bounds.Min.Y = point.Y
				}
				if point.Y > blobs[pos].Bounds.Max.Y {
					blobs[pos].Bounds.Max.Y = point.Y
				}
				blobs[pos].Points = append(blobs[pos].Points, point)
				blobs[pos].Centroid = blobs[pos].Bounds.Min.Add(blobs[pos].Bounds.Max).Div(2)
			}
		}
	}

	return &blobs
}

func (b *Blob) ClosestPoint(a image.Point) image.Point {
	b0 := b.Bounds

	//Get a rectangle containing point a
	b1 := image.Rectangle{
		Min: image.ZP,
		Max: a.Add(image.Pt(1, 1)),
	}

	//Get a rectangle contaning b0 and b1
	bounds := b0.Union(b1)

	//Generate a random color for each blob point
	nSites := len(b.Points)
	width, height := bounds.Dx(), bounds.Dy()
	colors := make([]color.NRGBA, nSites)
	rand.Seed(time.Now().Unix())
	for i := 0; i < nSites; i++ {
		colors[i] = color.NRGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)),
			uint8(rand.Intn(256)), 255}
	}

	//Generate a Voronoi Diagram by coloring each blob point with the color of the nearest blob point
	//Associate to each color a list of points
	m := map[color.NRGBA]([]image.Point){}
	img := image.NewNRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			dMin := dot(width, height)
			var colorMin int
			for s, p := range b.Points {
				if d := dot(p.X-x, p.Y-y); d < dMin {
					colorMin = s
					dMin = d
				}
			}
			img.SetNRGBA(x, y, colors[colorMin])
			if _, ok := m[colors[colorMin]]; !ok {
				m[colors[colorMin]] = []image.Point{b.Points[colorMin]}
			} else {
				m[colors[colorMin]] = append(m[colors[colorMin]], b.Points[colorMin])
			}
		}
	}

	minDist := math.MaxInt32
	nearestPoint := image.ZP

	//Get the point a color
	cellColor := img.NRGBAAt(a.X, a.Y)

	//Restore the list of points with this color
	for _, p := range m[cellColor] {

		//Find the nearest point to a
		dist := dot(p.X-a.X, p.Y-a.Y)
		if dist < minDist {
			minDist = dist
			nearestPoint = p
		}
	}

	return nearestPoint
}
