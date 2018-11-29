package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	gio "io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joaowiciuk/matrix"

	"github.com/joaowiciuk/vision"

	"github.com/anthonynsimon/bild/segment"

	"github.com/anthonynsimon/bild/effect"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/joaowiciuk/lenna/convolution"
	"github.com/joaowiciuk/lenna/morphology"
	o "github.com/joaowiciuk/lenna/threshold"
	lt "github.com/joaowiciuk/lenna/transform"
)

type params struct {
	a, b, c float64
}

func (p *params) String() string {
	return fmt.Sprintf("-lthres %.4f %.4f %.4f", p.a, p.b, p.c)
}

func (p *params) Set(values string) (err error) {
	aux := strings.Split(values, ",")
	if len(aux) != 3 {
		err = errors.New("Invalid params for local threshold")
	}
	p.a, _ = strconv.ParseFloat(aux[0], 64)
	p.b, _ = strconv.ParseFloat(aux[1], 64)
	p.c, _ = strconv.ParseFloat(aux[2], 64)
	/* fmt.Printf("%.4f %.4f %.4f\n", p.a, p.b, p.c) */
	if err != nil {
		err = errors.New("Invalid params for local threshold")
	}
	return
}

func main() {

	var kernel string
	var canny string
	var resize string
	var gray bool
	var threshold int
	var in string
	var out string
	var erode float64
	var dilate float64
	var invert bool
	var median float64
	var sobel bool
	var otsu bool
	var lthres params
	var conv string
	var morph string
	var ccl int
	var grad bool

	var width int
	var height int

	var flags = make(map[string]bool)
	var img image.Image
	var err error

	flag.StringVar(&kernel, "kernel", "data.csv", "-kernel data.csv")
	flag.StringVar(&canny, "canny", "91:31:3:1.4", "-canny hi:lo:w:s")
	flag.Var(&lthres, "lthres", "-lthres <stddev globmean locmean>")
	flag.BoolVar(&otsu, "otsu", false, "-otsu")
	flag.StringVar(&resize, "r", "683x384", "-r <comprimento>x<largura>")
	flag.BoolVar(&gray, "g", false, "-g")
	flag.IntVar(&threshold, "t", 0, "-t <corte>")
	flag.StringVar(&in, "in", "input.png", "-int <caminho/até/o/arquivo/nome_do_arquivo>.png")
	flag.StringVar(&out, "out", "output.png", "-out <output_name>.png")
	flag.Float64Var(&erode, "e", 1, "-e <valor>")
	flag.Float64Var(&dilate, "d", 1, "-d <valor>")
	flag.Float64Var(&median, "m", 1, "-m <valor>")
	flag.BoolVar(&invert, "i", false, "-i")
	flag.BoolVar(&sobel, "s", false, "-s")
	flag.StringVar(&conv, "conv", "boxblur", "-conv <kernel>")
	flag.StringVar(&morph, "morph", "erode,se.png", "-morph <operation> <kernel file>")
	flag.IntVar(&ccl, "ccl", 8, "-ccl <connectivity>")
	flag.BoolVar(&grad, "grad", false, "Image gradient magnitude")

	flag.Parse()

	width, _ = strconv.Atoi(strings.Split(resize, "x")[0])
	height, _ = strconv.Atoi(strings.Split(resize, "x")[1])

	flag.Visit(func(f *flag.Flag) { flags[f.Name] = true })

	if !flags["in"] || !flags["out"] {
		fmt.Printf("Arquivo de entrada ou saída não especificado\n")
		return
	}

	img, err = imgio.Open(in)

	if err != nil {
		fmt.Println("Erro ao abrir arquivo", err)
		return
	}

	if flags["r"] {
		img = transform.Resize(img, width, height, transform.Linear)
	}

	if flags["t"] {
		img = segment.Threshold(img, uint8(threshold))
	}

	if flags["g"] {
		img = lt.Grayscale(&img)
	}

	if flags["e"] {
		img = effect.Erode(img, erode)
	}

	if flags["d"] {
		img = effect.Dilate(img, dilate)
	}

	if flags["m"] {
		img = effect.Median(img, median)
	}

	if flags["i"] {
		img = effect.Invert(img)
	}

	if flags["s"] {
		img = effect.Sobel(img)
	}

	if flags["otsu"] {
		img = o.OtsuThreshold(&img)
	}

	if flags["lthres"] {
		grayimg := effect.Grayscale(img)
		img = o.Local(grayimg, o.Cumulative{Params: []float64{lthres.a, lthres.b, lthres.c}})
		/* fmt.Printf("%.4f %.4f %.4f\n", lthres.a, lthres.b, lthres.c) */
	}

	if flags["canny"] {
		params := strings.Split(canny, ":")
		if len(params) == 4 {
			hi, _ := strconv.Atoi(params[0])
			lo, _ := strconv.Atoi(params[1])
			n, _ := strconv.Atoi(params[2])
			s, _ := strconv.ParseFloat(params[3], 64)
			img = vision.Canny(img, uint8(hi), uint8(lo), n, s)
		}
	}

	if flags["kernel"] {
		f, _ := os.Open(kernel)
		r := csv.NewReader(f)
		lines := make([][]string, 0)
		m, n := 0, 0
		for {
			record, err := r.Read()
			if err == gio.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			m++
			if len(record) > n {
				n = len(record)
			}
			lines = append(lines, record)
		}
		x := matrix.New(m, n)
		for i0, v0 := range lines {
			for i1, v1 := range v0 {
				(*x)[i0][i1], _ = strconv.ParseFloat(v1, 64)
			}
		}
		/* fmt.Printf("%v\n", *x) */
		tensor := vision.Im2Mat(img)
		var src *matrix.Matrix
		switch len(tensor) {
		case 1:
			src = vision.Im2Mat(img)[0]
		case 3, 4:
			src.Law(func(r, c int) float64 {
				return 0.299*(*tensor[0])[r][c] + 0.587*(*tensor[1])[r][c] + 0.114*(*tensor[2])[r][c]
			})
		default:
			return
		}
		c := x.Conv(src)
		/* fmt.Printf("%v\n", *c) */
		img = *vision.Mat2Im([]*matrix.Matrix{c})
	}

	if flags["conv"] {
		var err error
		switch conv {
		case "laplacian":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.Laplacian())
		case "log":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.LaplacianOfGaussian())
		case "gblur33":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.GaussianBlur33())
		case "gblur55":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.GaussianBlur55())
		case "line135":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.Line135())
		case "sline":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.StraightLine())
		case "line45":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.Line45())
		case "hline":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.LineHorizontal())
		case "vline":
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.LineVertical())
		case "bblur":
			fallthrough
		default:
			img, err = convolution.Convolution(lt.Grayscale(&img), convolution.BoxBlur())
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if flags["morph"] {
		seFile := strings.Split(morph, ",")[1]
		se, _ := imgio.Open(seFile)
		sc := image.Pt(0, 0)
		if se.Bounds().Dx()%2 == 0 {
			sc.X = se.Bounds().Dx()/2 + 1
		} else {
			sc.X = se.Bounds().Dx() / 2
		}
		if se.Bounds().Dy()%2 == 0 {
			sc.Y = se.Bounds().Dy()/2 + 1
		} else {
			sc.Y = se.Bounds().Dy() / 2
		}
		switch strings.Split(morph, ",")[0] {
		case "open":
			img = morphology.Open(effect.Grayscale(img), lt.Grayscale(&se), sc)
		case "close":
			img = morphology.Close(effect.Grayscale(img), lt.Grayscale(&se), sc)
		case "dilate":
			img = morphology.Dilate(effect.Grayscale(img), lt.Grayscale(&se), sc)
		case "erode":
			fallthrough
		default:
			img = morphology.Erode(effect.Grayscale(img), lt.Grayscale(&se), sc)
		}
	}

	if flags["ccl"] {
		switch ccl {
		case 4:
			img = vision.Blobs(&img, vision.Connectivity4)
		default:
			img = vision.Blobs(&img, vision.Connectivity8)
		}
	}

	if flags["grad"] {
		gray := image.NewGray(img.Bounds())
		draw.Draw(gray, img.Bounds(), img, image.ZP, draw.Src)
		img, _ = vision.Grad(gray)
	}

	fileName := out[:strings.LastIndex(out, ".")]
	ext := out[strings.LastIndex(out, "."):]
	switch ext {
	case ".png":
		err = imgio.Save(fmt.Sprintf("%s.png", fileName), img, imgio.PNGEncoder())
	case ".jpg":
		err = imgio.Save(fmt.Sprintf("%s.jpg", fileName), img, imgio.JPEGEncoder(100))
	default:
		log.Println("Extensão de arquivo inválida")
	}
	if err != nil {
		log.Println("Erro ao salvar arquivo", err)
		return
	}
}
