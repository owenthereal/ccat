package main

import (
	"bufio"
	"fmt"
	"io"
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
	app.Action = runCCat

	app.Run(os.Args)
}

func runCCat(c *cli.Context) {
	var colorDefs ColorDefs
	if c.String("bg") == "dark" {
		colorDefs = DarkColorDefs
	} else {
		colorDefs = LightColorDefs
	}

	files := c.Args()
	// if there's no args, read from stdin
	if len(files) == 0 {
		files = append(files, "-")
	}

	for _, file := range files {
		err := ccat(file, colorDefs)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ccat(file string, colorDefs ColorDefs) error {
	var reader io.Reader
	if file == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		defer f.Close()
		reader = f
	}

	input, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	content, err := AsCCat(input, colorDefs)
	if err != nil {
		return err
	}

	fmt.Printf("%s", content)

	return nil
}
