package vision

import (
	"sort"

	"github.com/joaowiciuk/matrix"
)

type Region struct {
	D      []float64
	Xi, Xf []int
	Yi, Yf []int
}

func (R *Region) Len() int {
	return len((*R).D)
}

func (R *Region) Less(i, j int) bool {
	return (*R).D[i] < (*R).D[j]
}

func (R *Region) Swap(i, j int) {
	aux := Region{
		D:  make([]float64, 1),
		Xi: make([]int, 1),
		Xf: make([]int, 1),
		Yi: make([]int, 1),
		Yf: make([]int, 1),
	}
	aux.D[0] = (*R).D[i]
	aux.Xi[0] = (*R).Xi[i]
	aux.Xf[0] = (*R).Xf[i]
	aux.Yi[0] = (*R).Yi[i]
	aux.Yf[0] = (*R).Yf[i]

	(*R).D[i] = (*R).D[j]
	(*R).Xi[i] = (*R).Xi[j]
	(*R).Xf[i] = (*R).Xf[j]
	(*R).Yi[i] = (*R).Yi[j]
	(*R).Yf[i] = (*R).Yf[j]

	(*R).D[j] = aux.D[0]
	(*R).Xi[j] = aux.Xi[0]
	(*R).Xf[j] = aux.Xf[0]
	(*R).Yi[j] = aux.Yi[0]
	(*R).Yf[j] = aux.Yf[0]
}

/* func Find(O, S *matrix.Matrix) (R *Region) {

	xc, yc := matrix.Center(S)
	ms, ns := S.Size()
	mo, no := O.Size()

	left, right := xc, ns-xc-1
	top, bottom := yc, ms-yc-1

	R = &Region{
		D:  make([]float64, mo*no),
		Xi: make([]int, mo*no),
		Xf: make([]int, mo*no),
		Yi: make([]int, mo*no),
		Yf: make([]int, mo*no),
	}

	P := matrix.Padding(O, left, right, top, bottom)
	for x := left; x < no+left; x++ {
		for y := top; y < mo+top; y++ {
			S := matrix.SubMatrix(P, x-left, x+right+2, y-top, y+bottom+2)
			xi := x - left
			yi := y - top
			fmt.Printf("D[(%d,%d)-(%d,%d)] = %.4f\n", xi-left, xi+right+2, yi-top, yi+bottom+2, matrix.Distance(O, S, matrix.Norm1))
			(*R).D[yi*no+xi] = matrix.Distance(O, S, matrix.Norm1)
			(*R).Xi[yi*no+xi] = xi - left
			(*R).Xf[yi*no+xi] = xi + right + 2
			(*R).Yi[yi*no+xi] = yi - top
			(*R).Yf[yi*no+xi] = yi + bottom + 2
		}
	}
	for x := left; x < no+left; x++ {
		for y := top; y < mo+top; y++ {
			S := matrix.SubMatrix(P, x-left, x+right+2, y-top, y+bottom+2)
			d := matrix.Distance(O, S, matrix.Norm1)
			fmt.Println(d)
		}
	}

	return
} */

// Find finds O in I
/* func Find(O, I *matrix.Matrix) (R *Region) {
	xc, yc := matrix.Center(O)
	mo, no := O.Size()
	mi, ni := I.Size()

	left, right := xc, no-xc-1
	top, bottom := yc, mo-yc-1

	P := matrix.Padding(I, left, right, top, bottom)

	R = &Region{
		D:  make([]float64, mi*ni),
		Xi: make([]int, mi*ni),
		Xf: make([]int, mi*ni),
		Yi: make([]int, mi*ni),
		Yf: make([]int, mi*ni),
	}
	for x := left; x < ni+left; x++ {
		for y := top; y < mi+top; y++ {
			S := matrix.SubMatrix(P, x-left, x+right+2, y-top, y+bottom+2)
			d := matrix.Distance(S, O, matrix.Norm1)
			xi := x - left
			yi := y - top
			(*R).D[yi*ni+xi] = d
			(*R).Xi[yi*ni+xi] = xi - left
			(*R).Xf[yi*ni+xi] = xi + right + 2
			(*R).Yi[yi*ni+xi] = yi - top
			(*R).Yf[yi*ni+xi] = yi + bottom + 2
		}
	}

	return
} */

// Find finds O in I
func Find(O, I *matrix.Matrix, dist float64) (R *Region) {
	xc, yc := O.Center()
	mo, no := O.Size()
	mi, ni := I.Size()

	left, right := xc, no-xc-1
	top, bottom := yc, mo-yc-1

	P := I.Pad(left, right, top, bottom)

	R = &Region{
		D:  make([]float64, 0),
		Xi: make([]int, 0),
		Xf: make([]int, 0),
		Yi: make([]int, 0),
		Yf: make([]int, 0),
	}
	for x := left; x < ni+left; x++ {
		for y := top; y < mi+top; y++ {
			S := P.Submat(x-left, x+right+2, y-top, y+bottom+2)
			d := S.Dist(O, matrix.Norm1)
			xi := x - left
			yi := y - top
			if d <= dist {
				(*R).D = append((*R).D, d)
				(*R).Xi = append((*R).Xi, xi-left)
				(*R).Xf = append((*R).Xf, xi+right+2)
				(*R).Yi = append((*R).Yi, yi-top)
				(*R).Yf = append((*R).Yf, yi+bottom+2)
			}
		}
	}

	sort.Sort(R)
	return
}
