package pixelforging

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  // Para suporte a GIF
	_ "image/jpeg" // Para suporte a JPEG
	"image/png"
	_ "image/png" // Para suporte a PNG
	"log"
	"os"
	"sync"
)

// ListingPixels this function lists all the pixels in the image, this function does not preserve the order of the pixels
// if in the image the order of the pixels is X - Y in the final list of pixels it is possible that the pixels are shuffled and in the result the order may be Y - X
//
//	Use ListingPixelsOrded to get the list in the order of image
func ListingPixels(filePath string) ([]color.RGBA, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir a imagem:", err)
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Erro ao decodificar a imagem:", err)
		return nil, err
	}

	bounds := img.Bounds()
	var wg sync.WaitGroup
	colorsChan := make(chan color.RGBA, bounds.Dx()*10)
	sema := make(chan struct{}, 10)
	var colors []color.RGBA

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		sema <- struct{}{}
		go func(y int) {
			defer wg.Done()
			defer func() {
				<-sema
			}()

			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				colorsChan <- color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
			}

		}(y)

	}
	go func() {
		wg.Wait()
		close(colorsChan)
	}()

	for v := range colorsChan {
		colors = append(colors, v)
	}
	return colors, nil
}

// ListingPixelsOrded this function iterates through an image and returns all the pixels in a slice of color.RGBA, this function preserves the order in which the pixels appear in the image
func ListingPixelOrded(filePath string) ([]color.RGBA, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Erro ao abrir a imagem:", err)
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Erro ao decodificar a imagem:", err)
		return nil, err
	}

	bounds := img.Bounds()
	var colors []color.RGBA

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			colors = append(colors, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}

	}
	return colors, nil
}

func ExtractColorPalette(inputFilePath, outPutFilePath string) {
	colors, err := ListingPixels(inputFilePath)
	if err != nil {
		log.Fatalln("Error while trying to read the image pixels: ", err)
	}
	var uniqueColors []color.RGBA
	for _, color := range colors {
		addColorIfNotExists(&uniqueColors, color)
	}
	for _, color := range uniqueColors {
		fmt.Println(color.R, color.G, color.B, color.A)
	}
}

func addColorIfNotExists(uniqueColors *[]color.RGBA, color color.RGBA) {
	if len(*uniqueColors) == 0 {
		*uniqueColors = append(*uniqueColors, color)
		return
	} else {
		colorExists := false
		for _, v := range *uniqueColors {
			colorExists = v == color
		}
		if colorExists {
			return
		}
		*uniqueColors = append(*uniqueColors, color)
		return
	}
}

// CreateImage3x3 is a temp fuction for dev tests
func CreateImage3x3() error {
	// Cria uma nova imagem RGBA de tamanho 3x3.
	img := image.NewRGBA(image.Rect(0, 0, 3, 3))

	// Define as cores para cada linha.
	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}   // Vermelho
	green := color.RGBA{R: 0, G: 255, B: 0, A: 255} // Verde
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}  // Azul

	// Preenche os pixels da imagem.
	for x := 0; x < 3; x++ {
		img.Set(x, 0, red)   // Primeira linha: vermelho
		img.Set(x, 1, green) // Segunda linha: verde
		img.Set(x, 2, blue)  // Terceira linha: azul
	}

	// Cria o arquivo para salvar a imagem.
	file, err := os.Create("nova_imagem.png")
	if err != nil {
		return err
	}
	defer file.Close()

	// Salva a imagem no formato PNG.
	if err := png.Encode(file, img); err != nil {
		return err
	}

	return nil
}
