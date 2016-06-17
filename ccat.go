package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/mattn/go-isatty"
	"fmt"
)

type CCatPrinter interface {
	Print(r io.Reader, w io.Writer) error
}

type AutoColorPrinter struct {
	ColorPalettes ColorPalettes
}

func (a AutoColorPrinter) Print(r io.Reader, w io.Writer) error {
	if isatty.IsTerminal(uintptr(syscall.Stdout)) {
		return ColorPrinter{a.ColorPalettes}.Print(r, w)
	} else {
		return PlainTextPrinter{}.Print(r, w)
	}
}

type ColorPrinter struct {
	ColorPalettes ColorPalettes
}

func (c ColorPrinter) Print(r io.Reader, w io.Writer) error {
	return CPrint(r, w, c.ColorPalettes)
}

type PlainTextPrinter struct {
}

func (p PlainTextPrinter) Print(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}

type HtmlPrinter struct {
	ColorPalettes ColorPalettes
}

func (c HtmlPrinter) Print(r io.Reader, w io.Writer) error {
	return HtmlPrint(r, w, c.ColorPalettes)
}

func CCat(fname string, p CCatPrinter, w io.Writer) error {
	var r io.Reader

	if fname == readFromStdin {
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

		fi, err := file.Stat()

		if fi.Mode().IsDir() {
			return fmt.Errorf("%s is a directory", file.Name())
		}

		r = file
	}

	return p.Print(r, w)
}
