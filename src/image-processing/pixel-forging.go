package pixelforging

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  // Para suporte a GIF
	_ "image/jpeg" // Para suporte a JPEG
	"image/png"
	"log"
	"os"
	"sync"
)

type lineProcessingUtil struct {
	img    image.Image
	bounds image.Rectangle
}

// ListingPixels lists all the pixels in an image as a slice of `color.RGBA`.
// This function uses a worker pool to process the image in parallel, so the order
// of pixels in the result may not correspond to the original (row-column) order in the image.
// Use ListingPixelsOrdered if you need the pixel order to match the image's structure.
func ListingPixels(filePath string) ([]color.RGBA, error) {

	img, err := decodeImage(filePath)
	if err != nil {
		fmt.Println("Error while trying to open the image", err)
		return nil, err
	}
	var wg sync.WaitGroup
	bounds := img.Bounds()
	colorsChan := make(chan color.RGBA)
	colors := make([]color.RGBA, 0, bounds.Dx()*bounds.Dy())

	pools := 0
	numbersOfLines := bounds.Dy()

	if numbersOfLines < 32 {
		pools = numbersOfLines
	} else {
		pools = 32
	}

	linesToProcess := make(chan int, pools)

	for i := 0; i < pools; i++ {
		wg.Add(1)
		go woker(linesToProcess, colorsChan, lineProcessingUtil{
			img:    img,
			bounds: bounds,
		}, &wg)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		linesToProcess <- y
	}
	close(linesToProcess)

	go func() {
		wg.Wait()
		close(colorsChan)
	}()

	for v := range colorsChan {
		colors = append(colors, v)
	}
	return colors, nil
}

func woker(linesChan chan int, colorsChan chan color.RGBA, util lineProcessingUtil, wg *sync.WaitGroup) {
	defer wg.Done()
	for y := range linesChan {
		getPixelsOfImageLine(util.img, y, util.bounds, colorsChan)
	}

}

func getPixelsOfImageLine(img image.Image, line int, bounds image.Rectangle, colorsChan chan color.RGBA) {
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		r, g, b, a := img.At(x, line).RGBA()
		colorsChan <- color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}
}

// ListingPixelsOrdered iterates through an image and returns all its pixels as a slice of `color.RGBA`.
// The function preserves the original order of the pixels in the image (row by row).
func ListingPixelsOrdered(filePath string) ([]color.RGBA, error) {

	img, err := decodeImage(filePath)
	if err != nil {
		fmt.Println("Error while trying to open the image", err)
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

// ExtractColorPalette extracts the color palette of an image and prints the RGBA values of unique colors.
// Parameters:
// - inputFilePath: the path to the input image file.
// - outPutFilePath: the path to the output file (currently unused).
func ExtractColorPalette(inputFilePath, outPutFilePath string) {
	colors, err := ListingPixels(inputFilePath)
	if err != nil {
		log.Fatalln("Error while trying to read the image pixels: ", err)
	}

	uniqueColors := getUniqueColors(colors)
	for _, color := range uniqueColors {
		fmt.Println(color.R, color.G, color.B, color.A)
	}
}

// addColorIfNotExists adds a color to the slice of unique colors if it does not already exist.
func getUniqueColors(colors []color.RGBA) []color.RGBA {
	uniqueColors := make(map[color.RGBA]struct{})
	for _, color := range colors {
		uniqueColors[color] = struct{}{}
	}
	result := make([]color.RGBA, 0, len(uniqueColors))

	for key := range uniqueColors {
		result = append(result, key)
	}
	return result
}

func decodeImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir a imagem: %w", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar a imagem: %w", err)
	}
	return img, nil
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
