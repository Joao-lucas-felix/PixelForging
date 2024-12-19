package pixelforging

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  // Para suporte a GIF
	_ "image/jpeg" // Para suporte a JPEG
	_ "image/png"  // Para suporte a PNG
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
//ListingPixelsOrded this function iterates through an image and returns all the pixels in a slice of color.RGBA, this function preserves the order in which the pixels appear in the image
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
