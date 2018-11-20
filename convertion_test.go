package vision_test

import (
	"image"
	"testing"

	"github.com/joaowiciuk/matrix"
	"github.com/joaowiciuk/vision"
)

func TestIm2Mat(t *testing.T) {
	cases := []struct {
		desc     string
		value    image.Image
		expected []*matrix.Matrix
	}{
		{
			desc: "Gray",
			value: &image.Gray{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Pix: []uint8{
					0x80,
					0x80,
				},
			},
			expected: []*matrix.Matrix{
				&matrix.Matrix{
					{128.},
					{128.},
				},
			},
		},
		{
			desc: "Gray",
			value: &image.Gray{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Pix: []uint8{
					0x56,
					0x42,
				},
			},
			expected: []*matrix.Matrix{
				&matrix.Matrix{
					{86.},
					{66.},
				},
			},
		},
	}
	for _, c := range cases {
		actual := vision.Im2Mat(c.value)
		if len(actual) != len(c.expected) || !actual[0].Eq(c.expected[0]) {
			t.Errorf("%s: expected: %v, actual: %v", "Matrix from "+c.desc, c.expected, actual)
		}
	}
}

func TestMat2Im(t *testing.T) {
	cases := []struct {
		desc     string
		value    []*matrix.Matrix
		expected image.Image
	}{
		{
			desc: "1: 2x1",
			value: []*matrix.Matrix{
				&matrix.Matrix{
					{128.},
					{128.},
				},
			},
			expected: &image.Gray{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Pix: []uint8{
					0x80,
					0x80,
				},
			},
		},
		{
			desc: "2: 2x1",
			value: []*matrix.Matrix{
				&matrix.Matrix{
					{86.},
					{66.},
				},
			},
			expected: &image.Gray{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Pix: []uint8{
					0x56,
					0x42,
				},
			},
		},
	}
	for _, c := range cases {
		actual := vision.Mat2Im(c.value)
		if (*actual).Bounds().Dx() != c.expected.Bounds().Dx() || (*actual).Bounds().Dy() != c.expected.Bounds().Dy() {
			t.Errorf("%s: expected: %v, actual: %v", "Image from "+c.desc, c.expected, actual)
			continue
		}
		for x := 0; x < (*actual).Bounds().Dx(); x++ {
			for y := 0; y < (*actual).Bounds().Dy(); y++ {
				c1 := make([]uint32, 4)
				c1[0], c1[1], c1[2], c1[3] = (*actual).At(x, y).RGBA()
				c2 := make([]uint32, 4)
				c2[0], c2[1], c2[2], c2[3] = c.expected.At(x, y).RGBA()
				for i := 0; i < 4; i++ {
					if c1[i] != c2[i] {
						t.Errorf("%s: expected: %v, actual: %v", "Image from "+c.desc, c.expected, actual)
					}
				}
			}
		}
	}
}
