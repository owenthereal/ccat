package main

import (
	"log"
	"os"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/mattn/go-colorable"
)

const readFromStdin = "-"

var stdout = colorable.NewColorableStdout()

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

	fnames := c.Args()
	// if there's no args, read from stdin
	if len(fnames) == 0 {
		fnames = []string{readFromStdin}
	}

	for _, fname := range fnames {
		err := CCat(fname, colorDefs)
		if err != nil {
			log.Fatal(err)
		}
	}
}
