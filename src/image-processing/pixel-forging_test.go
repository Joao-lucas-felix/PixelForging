package pixelforging

import (
	"image"
	"image/color"
	"math/rand"
)

// Imagem 3x3 com 3 cores distintas.
// Imagem 10x10 com uma única cor.
// Imagem 10x10 com cores aleatórias e repetidas.
// Imagem transparente (10x10).
// Imagem maior (100x100) com gradiente de cores.

var (
	img3x3RGB          = image.NewRGBA(image.Rect(0, 0, 3, 3))
	img10x10Red        = image.NewRGBA(image.Rect(0, 0, 10, 10))
	img10x10A          = image.NewRGBA(image.Rect(0, 0, 10, 10))
	img10x10Rand       = image.NewRGBA(image.Rect(0, 0, 10, 10))
	img100x100Gradient = image.NewRGBA(image.Rect(0, 0, 100, 100))

	distinctColors3x3RGB          = []color.RGBA{red, green, blue}
	distinctColors10x10Red        = []color.RGBA{red}
	distinctColors10x10A          = []color.RGBA{transp}
	distinctColors10x10Rand       = make([]color.RGBA, 0, 100)
	distinctColors100x100Gradient = make([]color.RGBA, 0, 10.000)

	red    = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green  = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	blue   = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	transp = color.RGBA{R: 0, G: 0, B: 0, A: 0}
)

func CreateMock3x3RGB() {
	for x := 0; x < 3; x++ {
		img3x3RGB.Set(x, 0, red)
		img3x3RGB.Set(x, 1, green)
		img3x3RGB.Set(x, 2, blue)
	}
}
func CreateMock10x10RGB() {
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img10x10Red.Set(x, y, red)
		}
	}
}
func CreateMock10x10A() {
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img10x10A.Set(x, y, transp)
		}
	}
}
func CreateMock10x10Rand() {
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			randColor := color.RGBA{
				R: uint8(rand.Intn(256)),
				G: uint8(rand.Intn(256)),
				B: uint8(rand.Intn(256)),
				A: uint8(rand.Intn(256)),
			}
			img10x10Rand.Set(x, y, randColor)
			distinctColors10x10Rand = append(distinctColors10x10Rand, randColor)
		}
	}
}

func CreateMock100x100Gradient() {
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			colorGradient := color.RGBA{
				R: uint8(x),
				G: uint8(y),
				B: uint8((x + y) / 2),
				A: 255,
			}

			img100x100Gradient.Set(x, y, colorGradient)

			distinctColors100x100Gradient = append(distinctColors100x100Gradient, colorGradient)
		}
	}
}
