package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/mattn/go-colorable"
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/mattn/go-isatty"
)

var (
	stdout = colorable.NewColorableStdout()
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
		fnames = []string{"-"}
	}

	for _, fname := range fnames {
		err := ccat(fname, colorDefs)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ccat(fname string, colorDefs ColorDefs) error {
	var r io.Reader

	if fname == "-" {
		// scanner.Scanner from text/scanner couldn't detect EOF
		// if the io.Reader is os.Stdin
		// see https://github.com/golang/go/issues/10735
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		r = bytes.NewReader(b)
	} else {
		file, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer file.Close()
		r = file
	}

	r = bufio.NewReader(r)
	var err error
	if isatty.IsTerminal(uintptr(syscall.Stdout)) {
		var buf bytes.Buffer
		err = AsCCat(r, &buf, colorDefs)
		if err != nil {
			return err
		}

		_, err = stdout.Write(buf.Bytes())
	} else {
		_, err = io.Copy(os.Stdout, r)
	}

	return err
}
