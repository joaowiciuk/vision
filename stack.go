package vision

import "image"

func Stack(I image.RGBA, R image.Rectangle) (O []image.Image) {
	m, n := 0, 0
	w, h := R.Dx(), R.Dy()
	if I.Bounds().Dy()%h == 0 {
		m = int(I.Bounds().Dy() / h)
	} else {
		m = int(I.Bounds().Dy()/h) + 1
	}
	if I.Bounds().Dx()%w == 0 {
		n = int(I.Bounds().Dx() / w)
	} else {
		n = int(I.Bounds().Dx()/w) + 1
	}
	O = make([]image.Image, m*n)
	var frame, intersec image.Rectangle
	for c := 0; c < n; c++ {
		for r := 0; r < m; r++ {
			frame = image.Rect(c*w, r*h, (c+1)*w, (r+1)*h)
			intersec = I.Bounds().Intersect(frame)
			if !intersec.Empty() {
				if frame.Max.X > I.Bounds().Max.X {
					frame = frame.Sub(image.Pt(frame.Dx()-intersec.Dx(), 0))
				}
				if frame.Max.Y > I.Bounds().Max.Y {
					frame = frame.Sub(image.Pt(0, frame.Dy()-intersec.Dy()))
				}
			}
			O[r*n+c] = I.SubImage(frame)

		}
	}
	return
}
