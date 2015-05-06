package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ccat"
	app.Usage = "Concatenate FILE(s), or standard input, to standard output with colorized output."
	app.Version = Version
	app.Author = ""
	app.Email = ""
	app.HideHelp = true
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bg",
			Value: "light",
			Usage: `Set to "light" or "dark" depending on the terminal's background.`,
		},
	}
	app.Action = func(c *cli.Context) {
		var colorDefs ColorDefs
		if c.String("bg") == "dark" {
			colorDefs = DarkColorDefs
		} else {
			colorDefs = LightColorDefs
		}

		for _, file := range c.Args() {
			input, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}

			content, err := AsCCat(input, colorDefs)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s", content)
		}
	}

	app.Run(os.Args)
}
