package pixelforging

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"  // GIF
	"image/jpeg" // JPEG
	"image/png"  // PNG
	"io"
	"log"
	"math"
	"os"
	"sort"
	"sync"

	bmp "golang.org/x/image/bmp"   // BMP
	tiff "golang.org/x/image/tiff" // TIFF
	// WebP
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
	colorBlockWidth     = 50
	colorBlockHeight    = 50
	colorsPerRowDefault = 3
)

// ListingPixels lists all the pixels in an image as a slice of `color.RGBA`.
// This function uses a worker pool to process the image in parallel, so the order
// of pixels in the result may not correspond to the original (row-column) order in the image.
// Use ListingPixelsOrdered if you need the pixel order to match the image's structure.
func ListingPixels(img image.Image) ([]color.RGBA, error) {

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
		go worker(linesToProcess, colorsChan, lineProcessingUtil{
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

func worker(linesChan chan int, colorsChan chan color.RGBA, util lineProcessingUtil, wg *sync.WaitGroup) {
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

	img, err := DecodeImage(filePath)
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
func ExtractColorPalette(image image.Image, colorsPerRow, colorWidth, colorHeight int) image.Image {
	colors, err := ListingPixels(image)
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

	if image, err = createColorPalette(organizedColors, colorsPerRow, colorWidth, colorHeight); err != nil {
		log.Fatalln("Error wile trying to create the Collor Pallete", err)
	}
	return image
}

func createColorPalette(uniqueColors []color.RGBA, colorsPerRow, colorWidth, colorHeight int) (image.Image, error) {

	colorBlocks := make([]image.Image, 0, len(uniqueColors))
	// Creates c blocks to create the palette

	// refactor idea:
	// i can use the worker pools pattern here
	// 1 chanel of uniqueColors
	// the worker create the c block and add to a chanel of image
	var wg sync.WaitGroup
	colorsChan := make(chan color.RGBA, len(uniqueColors))
	imagesChan := make(chan image.Image, len(uniqueColors))
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(colorsChan chan color.RGBA) {
			defer wg.Done()
			for c := range colorsChan {
				img := image.NewRGBA(image.Rect(0, 0, colorWidth, colorHeight))
				for y := 0; y < img.Bounds().Dy(); y++ {
					for x := 0; x < img.Bounds().Dx(); x++ {
						img.Set(x, y, c)
					}
				}
				imagesChan <- img
			}
		}(colorsChan)
	}

	for _, c := range uniqueColors {
		colorsChan <- c
	}
	close(colorsChan)

	go func() {
		for img := range imagesChan {
			colorBlocks = append(colorBlocks, img)
		}
		fmt.Println("Parallel processing complete")
	}()
	wg.Wait()
	close(imagesChan)

	var verticalColors []image.Image
	horizontalColors := make([]image.Image, 0, len(colorBlocks)/4)
	// {} {} {} {} {}
	for _, c := range colorBlocks {
		verticalColors = append(verticalColors, c)
		if len(verticalColors) == colorsPerRow {
			horizontalColor, err := concatenateImagesHorizontal(colorHeight, verticalColors...)
			if err != nil {
				return nil, err
			}
			horizontalColors = append(horizontalColors, horizontalColor)
			verticalColors = verticalColors[:0]
		}
	}
	// Process any remaining vertical colors
	if len(verticalColors) > 0 {
		horizontalColor, err := concatenateImagesHorizontal(colorHeight, verticalColors...)
		if err != nil {
			return nil, err

		}
		horizontalColors = append(horizontalColors, horizontalColor)
	}

	img, err := concatenateImagesVertical(colorWidth, colorsPerRow, horizontalColors...)
	if err != nil {
		return nil, err

	}
	return img, nil
}

func concatenateImagesHorizontal(colorHeight int, imgs ...image.Image) (image.Image, error) {
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

func concatenateImagesVertical(colorWidth, colorsPerRow int, imgs ...image.Image) (image.Image, error) {
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

// SaveImage save the image on Output file path
func SaveImage(img image.Image, outPutFilePath string) error {
	file, err := os.Create(outPutFilePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("Error while trying to close file", err)
		}
	}(file)

	// Salva a imagem no formato PNG.
	if err := png.Encode(file, img); err != nil {
		return err
	}

	return nil
}

// addColorIfNotExists adds a color to the slice of unique colors if it does not already exist.
func getUniqueColors(colors []color.RGBA) []color.RGBA {
	uniqueColors := make(map[color.RGBA]struct{})
	for _, c := range colors {
		uniqueColors[c] = struct{}{}
	}
	result := make([]color.RGBA, 0, len(uniqueColors))

	for key := range uniqueColors {
		result = append(result, key)
	}
	return result
}

// DecodeImage open a image from a path
func DecodeImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir a imagem: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("Error while trying to close file", err)
		}
	}(file)
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("Error while trying to close file", err)
		}
	}(file)

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

	maxRGBA := math.Max(r, math.Max(g, b))
	minRGBA := math.Min(r, math.Min(g, b))
	l = (maxRGBA + minRGBA) / 2

	if maxRGBA == minRGBA {
		h, s = 0, 0
	} else {
		delta := maxRGBA - minRGBA
		s = delta / (1 - math.Abs(2*l-1))
		switch maxRGBA {
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
		if c.R != 0 && c.G != 0 && c.B != 0 && c.A != 0 {
			h, s, l := RGBAToHSL(c)
			hslColors[i] = HSLColor{Color: c, H: h, S: s, L: l}
		}
	}

	// Ordenar por tonalidade (H), saturação (S) e luminosidade (L)
	sort.Slice(hslColors, func(i, j int) bool {
		if hslColors[i].H == hslColors[j].H {
			return hslColors[i].L < hslColors[j].L
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

// BytesToImage decodes an image from a byte slice.
// It detects the image format and decodes it accordingly.
// It returns the decoded image, its format, and any error encountered.
// The function supports JPEG, PNG, GIF, BMP, TIFF, and WebP formats.
// If the format is not recognized, it attempts a generic decode.
// The function uses a bytes.Reader to read the image data from the byte slice.
// It returns an error if the image cannot be decoded or if the format is not supported.
func BytesToImage(imgBytes []byte,fm string) (image.Image, string, error) {
	imgReader := bytes.NewReader(imgBytes)

	// Tenta detectar o formato
	_, format, err := image.DecodeConfig(imgReader)
	if err != nil {
		return nil, "", err
	}

	// Volta ao início do reader
	imgReader.Seek(0, io.SeekStart)

	// Decodifica usando o decoder específico para melhor controle
	switch format {
	case "jpeg":
		img, err := jpeg.Decode(imgReader)
		return img, "", err
	case "png":
		img, err := png.Decode(imgReader)
		return img, "", err
	case "gif":
		img, err := gif.Decode(imgReader)
		return img, "", err
	case "bmp":
		img, err := bmp.Decode(imgReader)
		return img, "", err
	case "tiff":
		img, err := tiff.Decode(imgReader)
		return img, "", err
	default:
		// Tenta decodificação genérica
		return image.Decode(imgReader)
	}
}

// ImageToBytes converte uma image.Image para bytes no formato especificado
// format: "jpeg", "png", "gif", "bmp", "tiff", "webp"
func ImageToBytes(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	switch format {
	case "jpeg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(&buf, img)
	case "gif":
		// O pacote image/gif tem funções mais complexas para GIFs animados
		// Esta é uma implementação básica para GIFs estáticos
		err = gif.Encode(&buf, img, &gif.Options{})
	case "bmp":
		err = bmp.Encode(&buf, img)
	case "tiff":
		err = tiff.Encode(&buf, img, &tiff.Options{})
	default:
		// Default para PNG se o formato não for reconhecido
		err = png.Encode(&buf, img)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
