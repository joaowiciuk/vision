package kernel

import (
	"math"
	"reflect"
	"testing"

	"github.com/joaowiciuk/matrix"
)

func TestGaussian(t *testing.T) {
	type args struct {
		n int
		σ float64
	}
	tests := []struct {
		name string
		args args
		want *matrix.Matrix
	}{
		{
			name: "3x3, σ = 1.0",
			args: args{
				n: 3,
				σ: 1.0,
			},
			want: &matrix.Matrix{
				{0.0751, 0.1238, 0.0751},
				{0.1238, 0.2042, 0.1238},
				{0.0751, 0.1238, 0.0751},
			},
		},
		{
			name: "5x5, σ: 1.0",
			args: args{
				n: 5,
				σ: 1.0,
			},
			want: &matrix.Matrix{
				{0.0030, 0.0133, 0.0219, 0.0133, 0.0030},
				{0.0133, 0.0596, 0.0983, 0.0596, 0.0133},
				{0.0219, 0.0983, 0.1621, 0.0983, 0.0219},
				{0.0133, 0.0596, 0.0983, 0.0596, 0.0133},
				{0.0030, 0.0133, 0.0219, 0.0133, 0.0030},
			},
		},
		{
			name: "3x3, σ: 1.4",
			args: args{
				n: 3,
				σ: 1.4,
			},
			want: &matrix.Matrix{
				{0.0924, 0.1192, 0.0924},
				{0.1192, 0.1538, 0.1192},
				{0.0924, 0.1192, 0.0924},
			},
		},
		{
			name: "3x3, σ: 1.7",
			args: args{
				n: 3,
				σ: 1.7,
			},
			want: &matrix.Matrix{
				{0.0983, 0.1169, 0.0983},
				{0.1169, 0.1390, 0.1169},
				{0.0983, 0.1169, 0.0983},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Gaussian(tt.args.n, tt.args.σ); !got.ErrEq(tt.want, 1e-3) {
				t.Errorf("Gaussian() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestLaplacian(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Laplacian(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Laplacian() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSharpen(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sharpen(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sharpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine180(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Line180(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line180() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine90(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Line90(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line90() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine45(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Line45(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line45() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine135(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Line135(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Line135() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoG(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoG(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoG() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBox(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Box(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Box() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsharp55(t *testing.T) {
	tests := []struct {
		name string
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unsharp55(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unsharp55() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSobelX(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SobelX(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SobelX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSobelY(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want *matrix.Matrix
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SobelY(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SobelY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGaussian1D(t *testing.T) {
	type args struct {
		n int
		σ float64
	}
	tests := []struct {
		name  string
		args  args
		wantA *matrix.Matrix
	}{
		{
			args: args{
				n: 7,
				σ: math.Sqrt(2),
			},
			wantA: &matrix.Matrix{
				{1, 4, 7, 10, 7, 4, 1},
				{4, 12, 26, 33, 26, 12, 4},
				{7, 26, 55, 71, 55, 26, 7},
				{10, 33, 71, 91, 71, 33, 10},
				{7, 26, 55, 71, 55, 26, 7},
				{4, 12, 26, 33, 26, 12, 4},
				{1, 4, 7, 10, 7, 4, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA := Gaussian1D(tt.args.n, tt.args.σ)
			if !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("Gaussian1D() gotA = \n%v\nWant \n%v\n", gotA, tt.wantA)
			}
		})
	}
}
