package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Joao-lucas-felix/PixelForging/src/backend/server"
	pixelforging "github.com/Joao-lucas-felix/PixelForging/src/image-processing"
	"github.com/urfave/cli"
)
const (
	color = "\033[38;5;117m"
	colorReset = "\033[0m"
	logo = color+`
██████╗ ██╗██╗  ██╗███████╗██╗     ███████╗ ██████╗ ██████╗  ██████╗ ██╗███╗   ██╗ ██████╗ 
██╔══██╗██║╚██╗██╔╝██╔════╝██║     ██╔════╝██╔═══██╗██╔══██╗██╔════╝ ██║████╗  ██║██╔════╝ 
██████╔╝██║ ╚███╔╝ █████╗  ██║     █████╗  ██║   ██║██████╔╝██║  ███╗██║██╔██╗ ██║██║  ███╗
██╔═══╝ ██║ ██╔██╗ ██╔══╝  ██║     ██╔══╝  ██║   ██║██╔══██╗██║   ██║██║██║╚██╗██║██║   ██║
██║     ██║██╔╝ ██╗███████╗███████╗██║     ╚██████╔╝██║  ██║╚██████╔╝██║██║ ╚████║╚██████╔╝
╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝╚═╝  ╚═══╝ ╚═════╝ 
																					
`+colorReset
)


// GenAPP gen a news cli app
func GenAPP() *cli.App {
	app := cli.NewApp()
	app.Name = "Pixel Forging"

	app.Usage = "A CLI to processing images"

	app.Commands = []cli.Command{
		//Hello comand
		{
			Name:  "hello",
			Usage: "Greets the user to check the safety of app, you can use the --name=\"[YOUR_NAME}\" to this methodo greets you! ",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "World",
				},
			},
			Action: func(c *cli.Context) {
				name := c.String("name")
				fmt.Println("Hello ", name)
			},
		},
		// Extract palette command
		{
			Name:  "extract-palette",
			Usage: "Opens the image in the dir that you pass in the flag --input-image=\"[YOUR-IMAGE_PATH}\" and extract the color palette of the image and saves in the path that you pass in the flag --output-image=\"[OUTPUT_IMAGE_PATH]\"\nYou can pass 3 parans to configure the size of palette color image:\n\t--colors-per-row=\"[NUMBER_OF_COLORS_PER_ROW]\"\n\t--width=\"[WIDTH_OF_COLOR_BLOCK]\"\n\t--height=\"[HEIGHT_OF_COLOR_BLOCK]\" \n  --colors-num=\"[NUMBER_OF_COLORS]\"\n\nThe default values are:\n\t--colors-per-row=3\n\t--width=0\n\t--height=0\n\t--colors-num=0",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input-image",
					Value: "",
				},
				cli.StringFlag{
					Name:  "output-image",
					Value: "",
				},
				cli.StringFlag{
					Name:  "colors-per-row",
					Value: "0",
				},
				cli.StringFlag{
					Name:  "width",
					Value: "0",
				},
				cli.StringFlag{
					Name:  "height",
					Value: "0",
				},
				cli.StringFlag{
					Name:  "colors-num",
					Value: "0",
				},
			},
			Action: func(c *cli.Context) {
				fmt.Println(logo)
				inputPath := c.String("input-image")
				outputPath := c.String("output-image")
				colorsPerRowS := c.String("colors-per-row")
				widthS := c.String("width")
				heightS := c.String("height")
				colorNumS := c.String("colors-num")

				if inputPath == "" {
					log.Fatalln("The param --input-image can not be blanck")
				}

				if outputPath == "" {
					log.Fatalln("The param --output-image can not be blanck")
				}

				colorsPerRow, err := strconv.Atoi(colorsPerRowS)
				if err != nil {
					log.Fatalln("The param --colors-per-row should be a int number")
				}

				width, err := strconv.Atoi(widthS)
				if err != nil {
					log.Fatalln("The param --width should be a int number")
				}
				height, err := strconv.Atoi(heightS)
				if err != nil {
					log.Fatalln("The param --height should be a int number")
				}
				colorNum, err := strconv.Atoi(colorNumS)
				if err != nil {
					log.Fatalln("The param --colors-num should be a int number")
				}
				image, err := pixelforging.DecodeImage(inputPath)
				if err != nil {
					log.Fatalln(err)
				}

				fmt.Println("We are forging your palette!")

				img := pixelforging.ExtractColorPalette(image, colorsPerRow, width, height, colorNum)

				if err := pixelforging.SaveImage(img, outputPath); err != nil {
					log.Fatalln(err)
				}
			},
		},
		// Init server command
		{
			Name:  "start-gRPC-server",
			Usage: "Starts the gRPC server on a port that you pass in the flag --port=\"[PORT]\" the default port is 9090",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Value: "9090",
				},
			},
			Action: func(c *cli.Context) {
				fmt.Println(logo)
				port := c.String("port")
				server.BoostrapServer(port)
			},
		},
	}
	return app
}
