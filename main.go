package main

import (
	"bufio"
	"log"
	"os"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/codegangsta/cli"
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

	fnames := c.Args()
	// if there's no args, read from stdin
	if len(fnames) == 0 {
		fnames = append(fnames, "-")
	}

	for _, fname := range fnames {
		err := ccat(fname, colorDefs)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ccat(fname string, colorDefs ColorDefs) error {
	file := os.Stdin
	if fname != "-" {
		var err error
		file, err = os.Open(fname)
		if err != nil {
			return err
		}

		defer file.Close()
	}

	content, err := AsCCat(bufio.NewReader(file), colorDefs)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(content)

	return err
}
