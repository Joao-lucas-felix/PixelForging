package app

import (
	"fmt"

	"github.com/urfave/cli"
)

// Gerar um novo app
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
	}
	return app
}
