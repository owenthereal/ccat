package main

import (
	"log"
	"os"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/mattn/go-colorable"
)

const (
	readFromStdin = "-"
)

func init() {
	cli.AppHelpTemplate = `NAME:
    {{.Name}} - {{.Usage}}

USAGE:
    {{.Name}} [options] [file ...]

VERSION:
    {{.Version}}

OPTIONS:
	{{range .Flags}}{{.}}
	{{end}}
Using color is auto both by default and with --color=auto. With --color=auto,
ccat emits color codes only when standard output is connected to a terminal.
`
}

func main() {
	app := cli.NewApp()
	app.Name = "ccat"
	app.Usage = "Colorize FILE(s), or standard input, to standard output."
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bg",
			Value: "light",
			Usage: `Set to "light" or "dark" depending on the terminal's background.`,
		},
		cli.StringFlag{
			Name:  "C, color",
			Value: "auto",
			Usage: `Colorize the output; Value can be "never", "always" or "auto".`,
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

	var printer CCatPrinter
	if c.String("color") == "always" {
		printer = ColorPrinter{colorDefs}
	} else if c.String("color") == "never" {
		printer = PlainTextPrinter{}
	} else {
		printer = AutoColorPrinter{colorDefs}
	}

	fnames := c.Args()
	// if there's no args, read from stdin
	if len(fnames) == 0 {
		fnames = []string{readFromStdin}
	}

	stdout := colorable.NewColorableStdout()
	for _, fname := range fnames {
		err := CCat(fname, printer, stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}
