package main

import (
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
	
	pixelforging.ExtractColorPalette("nova_imagem.png", "")
}
