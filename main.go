package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Joao-lucas-felix/PixelForging/src/app"
	pixelforging "github.com/Joao-lucas-felix/PixelForging/src/image-processing"
)

func main() {
	aplication := app.GenAPP()
	erro := aplication.Run(os.Args)
	if erro != nil {
		log.Fatal(erro)
	}

	clors, _ := pixelforging.ListingPixels("nova_imagem.png")
	for _, color := range clors {
		fmt.Println(color.R, color.G, color.B, color.A)
	}

	clors, _ = pixelforging.ListingPixelOrded("nova_imagem.png")
	for _, color := range clors {
		fmt.Println(color.R, color.G, color.B, color.A)
	}
}
