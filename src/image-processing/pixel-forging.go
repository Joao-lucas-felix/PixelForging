package pixelforging

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"  // Para suporte a GIF
	_ "image/jpeg" // Para suporte a JPEG
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"sync"
)

type lineProcessingUtil struct {
	img    image.Image
	bounds image.Rectangle
}

type HSLColor struct {
	Color   color.RGBA
	H, S, L float64
}

const (
	colorBlockWidth  = 50
	colorBlockHeight = 50
	colorsPerRowDefault    = 3
)

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
	colorsChan := make(chan color.RGBA, bounds.Dx()*bounds.Dy()/32)
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

	go func() {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			linesToProcess <- y
		}
		close(linesToProcess)
	}()

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
func ExtractColorPalette(inputFilePath, outPutFilePath string, colorsPerRow, colorWidth, colorHeight int) {
	colors, err := ListingPixels(inputFilePath)
	if err != nil {
		log.Fatalln("Error while trying to read the image pixels: ", err)
	}
	if colorsPerRow == 0 {
		colorsPerRow = colorsPerRowDefault
	}
	if colorWidth == 0 {
		colorWidth = colorBlockWidth
	}
	if colorHeight == 0 {
		colorHeight = colorBlockHeight
	}

	uniqueColors := getUniqueColors(colors)
	organizedColors := organizeColorsByHSL(uniqueColors)
	if err := createColorPalette(organizedColors, outPutFilePath, colorsPerRow, colorWidth, colorHeight); err != nil {
		log.Fatalln("Error wile trying to create the Collor Pallete", err)
	}

}

func createColorPalette(uniqueColors []color.RGBA, outPutFilePath string, colorsPerRow, colorWidth, colorHeight int) error {

	colorBlocks := make([]image.Image, 0, len(uniqueColors))
	// Creates color blocks to create the palette
	for _, color := range uniqueColors {
		img := image.NewRGBA(image.Rect(0, 0, colorWidth, colorHeight))
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				img.Set(x, y, color)
			}
		}
		colorBlocks = append(colorBlocks, img)
	}

	var verticalColors []image.Image
	horizontalColors := make([]image.Image, 0, len(colorBlocks)/4)
	// {} {} {} {} {}
	for _, color := range colorBlocks {
		verticalColors = append(verticalColors, color)
		if len(verticalColors) == colorsPerRow {
			horizontalColor, err := concatenateImagesHorizontal(colorHeight, verticalColors...)
			if err != nil {
				return err
			}
			horizontalColors = append(horizontalColors, horizontalColor)
			verticalColors = verticalColors[:0]
		}
	}
	// Process any remaining vertical colors
	if len(verticalColors) > 0 {
		horizontalColor, err := concatenateImagesHorizontal(colorHeight, verticalColors...)
		if err != nil {
			return err
		}
		horizontalColors = append(horizontalColors, horizontalColor)
	}

	img, err := concatenateImagesVertical(colorWidth, colorsPerRow, horizontalColors...)
	if err != nil {
		return err
	}
	if err := saveImage(img, outPutFilePath); err != nil {
		return err
	}
	return nil
}

func concatenateImagesHorizontal(colorHeight int ,imgs ...image.Image) (image.Image, error) {
	if len(imgs) == 0 {
		return nil, fmt.Errorf("no images to concatenate")
	}
	for _, img := range imgs {
		if img == nil {
			return nil, fmt.Errorf("nil image in input in horizontal")
		}
	}

	// Define a nova largura e altura para a imagem concatenada
	width := 0
	for _, img := range imgs {
		width += img.Bounds().Dx()
	}

	// Cria uma nova imagem vazia
	newImage := image.NewRGBA(image.Rect(0, 0, width, colorHeight))

	offsetX := 0
	for _, img := range imgs {
		draw.Draw(newImage, img.Bounds().Add(image.Pt(offsetX, 0)), img, image.Point{}, draw.Src)
		offsetX += img.Bounds().Dx()
	}

	return newImage, nil
}

func concatenateImagesVertical(colorWidth, colorsPerRow int ,imgs ...image.Image) (image.Image, error) {
	if len(imgs) == 0 {
		return nil, fmt.Errorf("no images to concatenate")
	}
	for _, img := range imgs {
		if img == nil {
			return nil, fmt.Errorf("nil image in input in vertical")
		}
	}

	// Define a nova largura e altura para a imagem concatenada
	height := 0
	for _, img := range imgs {
		height += img.Bounds().Dy()
	}

	// Cria uma nova imagem vazia
	newImage := image.NewRGBA(image.Rect(0, 0, colorWidth*colorsPerRow, height))

	// Desenha a primeira imagem

	offsetY := 0
	// Desenha a segunda imagem, deslocada para baixo
	for _, img := range imgs {
		draw.Draw(newImage, img.Bounds().Add(image.Pt(0, offsetY)), img, image.Point{}, draw.Src)
		offsetY += img.Bounds().Dy()
	}
	return newImage, nil
}

func saveImage(img image.Image, outPutFilePath string) error {
	file, err := os.Create(outPutFilePath)
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
// RGBAToHSL converts an color in RGBA space to HSL 
func RGBAToHSL(c color.RGBA) (h, s, l float64) {
	r := float64(c.R) / 255
	g := float64(c.G) / 255
	b := float64(c.B) / 255

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l = (max + min) / 2

	if max == min {
		h, s = 0, 0 
	} else {
		delta := max - min
		s = delta / (1 - math.Abs(2*l-1))
		switch max {
		case r:
			h = math.Mod((g-b)/delta+6, 6)
		case g:
			h = (b-r)/delta + 2
		case b:
			h = (r-g)/delta + 4
		}
		h *= 60 
	}

	return
}

func organizeColorsByHSL(colors []color.RGBA) []color.RGBA {
	hslColors := make([]HSLColor, len(colors))
	for i, c := range colors {
		h, s, l := RGBAToHSL(c)
		hslColors[i] = HSLColor{Color: c, H: h, S: s, L: l}
	}

	// Ordenar por tonalidade (H), saturação (S) e luminosidade (L)
	sort.Slice(hslColors, func(i, j int) bool {
		if hslColors[i].L != hslColors[j].L {
			return hslColors[i].L < hslColors[j].L
		}
		if hslColors[i].S != hslColors[j].S {
			return hslColors[i].S < hslColors[j].S
		}
		return hslColors[i].H < hslColors[j].H
	})

	// Retorna a lista ordenada
	sortedColors := make([]color.RGBA, len(colors))
	for i, c := range hslColors {
		sortedColors[i] = c.Color
	}
	return sortedColors
}
