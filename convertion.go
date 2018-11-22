package vision

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/anthonynsimon/bild/math/f64"
	"github.com/joaowiciuk/matrix"
)

// Im2Mat converts an image to an array of matrices.
// The output array contains a single matrix if the input is a grayscale image
// or four matrices if it is a RGBA image.
// When converting from RGBA, the channel to index correspondence is the following:
//		RED 	-> 0
//		GREEN 	-> 1
//		BLUE 	-> 2
//		ALPHA 	-> 3
// If the image pointer is nil the function produces nothing and returns a nil pointer, so remember always
// checking for nil when using it.
func Im2Mat(i image.Image) (array []*matrix.Matrix) {
	if i == nil {
		return
	}
	b := i.Bounds()
	switch i.ColorModel() {
	case color.GrayModel, color.Gray16Model:
		array = make([]*matrix.Matrix, 1)
		array[0] = matrix.New(b.Dy(), b.Dx())
		imgGRAY := image.NewGray(b)
		draw.Draw(imgGRAY, b, i, b.Min, draw.Src)
		for x := 0; x < b.Dx(); x++ {
			for y := 0; y < b.Dy(); y++ {
				c := imgGRAY.GrayAt(x, y)
				(*array[0])[y][x] = float64(c.Y)
			}
		}
	default:
		array = make([]*matrix.Matrix, 4)
		for i := range array {
			array[i] = matrix.New(b.Dy(), b.Dx())
		}
		imgRGBA := image.NewRGBA(b)
		draw.Draw(imgRGBA, b, i, b.Min, draw.Src)
		for x := 0; x < b.Dx(); x++ {
			for y := 0; y < b.Dy(); y++ {
				C := imgRGBA.RGBAAt(x, y)
				(*array[0])[y][x] = float64(C.R) //RED
				(*array[1])[y][x] = float64(C.G) //GREEN
				(*array[2])[y][x] = float64(C.B) //BLUE
				(*array[3])[y][x] = float64(C.A) //ALPHA
			}
		}
	}
	return
}

// Gray2Mat converts a gray scale image to matrix.
func Gray2Mat(i *image.Gray) *matrix.Matrix {
	if i == nil {
		return nil
	}
	b := i.Bounds()
	mat := matrix.New(b.Dy(), b.Dx())
	for x := 0; x < b.Dx(); x++ {
		for y := 0; y < b.Dy(); y++ {
			c := i.GrayAt(x, y)
			mat.Set(y, x, float64(c.Y))
		}
	}
	return mat
}

// Mat2Gray converts a matrix into a gray scale image.
func Mat2Gray(mat *matrix.Matrix) (gray *image.Gray) {
	rows, cols := mat.Size()
	gray = image.NewGray(image.Rect(0, 0, cols, rows))
	for i := range gray.Pix {
		row, col := mat.ToPoint(i)
		gray.Pix[i] = uint8(clamp(mat.At(row, col), 0., 255.)) //GRAY
	}
	return
}

// Mat2Im converts an valid array of matrices to an 8-bit color depth image and returns a pointer to it.
// A valid array contains matrices of same size.
// The function produces a grayscale image if the array has a single matrix and
// a RGBA image if it has four matrices.
// When converting to RGBA, the index to channel correspondence is the following:
//		0 -> RED
//		1 -> GREEN
//		2 -> BLUE
//		3 -> ALPHA
// If the array is not valid the function produces nothing and returns a nil pointer, so remember always
// checking for nil when using it.
func Mat2Im(array []*matrix.Matrix) (ptr *image.Image) {
	c := len(array)

	//Checks if the input array is valid in lenght
	if c != 1 && c != 4 {
		ptr = nil
		return
	}

	//Checks if the input array is valid in elements dimensions
	m, n := make([]int, c), make([]int, c)
	for i, v := range array {
		m[i], n[i] = v.Size()
		if i == 0 {
			continue
		}
		if m[i] != m[i-1] || n[i] != n[i-1] {
			ptr = nil
			return
		}
	}

	//Since here its checked that all elements have same dimension, any element
	//can be used to calculate the produced image bounds. Below the last
	//element is used. Note the difference between rows and columns in an
	//matrix and (x, y) coordinates in an image
	b := image.Rect(0, 0, n[c-1], m[c-1])

	switch c {
	case 4:
		imgRGBA := image.NewRGBA(b)
		n := 0
		for i := range imgRGBA.Pix {
			if i%4 == 0 {
				//Here a for loop could be used, but for lucidity purposes the values
				//are written explicitly. Ctrl + c, Ctrl + v helped a lot
				//A matrix can pass for some operations like Laplacian, Gaussian, etc,
				//thus it values should be numericaly "rearranged" for matching a 8-bit depth color,
				//thus the use of Clamp function
				row, col := array[0].ToPoint(n)
				imgRGBA.Pix[i+0] = uint8(f64.Clamp(array[0].At(row, col), 0., 255.)) //RED
				imgRGBA.Pix[i+1] = uint8(f64.Clamp(array[1].At(row, col), 0., 255.)) //GREEN
				imgRGBA.Pix[i+2] = uint8(f64.Clamp(array[2].At(row, col), 0., 255.)) //BLUE
				imgRGBA.Pix[i+3] = uint8(f64.Clamp(array[3].At(row, col), 0., 255.)) //ALPHA
				n++
			}
		}
		aux := image.Image(imgRGBA)
		ptr = &aux
	case 1:
		imgGray := image.NewGray(b)
		for i := range imgGray.Pix {
			row, col := array[0].ToPoint(i)
			imgGray.Pix[i] = uint8(f64.Clamp(array[0].At(row, col), 0., 255.)) //GRAY
		}
		aux := image.Image(imgGray)
		ptr = &aux
	}
	return
}
