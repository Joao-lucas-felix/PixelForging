package main

import (
	"log"
	"os"

	"github.com/Joao-lucas-felix/PixelForging/src/app"
)

func main() {
	aplication := app.GenAPP()
	erro := aplication.Run(os.Args)
	if erro != nil {
		log.Fatal(erro)
	}
}
