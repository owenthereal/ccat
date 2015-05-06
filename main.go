package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"syscall"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/jingweno/ccat/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"
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
	file := os.Stdin
	if fname != "-" {
		var err error
		file, err = os.Open(fname)
		if err != nil {
			return err
		}

		defer file.Close()
	}

	reader := bufio.NewReader(file)

	var err error
	if terminal.IsTerminal(syscall.Stdout) {
		var buf bytes.Buffer
		err = AsCCat(reader, &buf, colorDefs)
		if err != nil {
			return err
		}

		_, err = os.Stdout.Write(buf.Bytes())
	} else {
		_, err = io.Copy(os.Stdout, reader)
	}

	return err
}
