package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"syscall"
	"fmt"

	"github.com/mattn/go-isatty"
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

func CCat(fname string, p CCatPrinter, w io.Writer, s int) error {
	var r io.Reader

	if fname == readFromStdin {
		b, err := ioutil.ReadAll(os.Stdin)
		if s > 0 && len(b) > s {
			return errors.New(fmt.Sprintf("read size is greater than %d bytes\n", s))
		}
		if err != nil {
			return err
		}

		r = bytes.NewReader(b)
	} else {
		if fstat, _ := os.Stat(fname); s > 0 && fstat.Size() > int64(s) {
			return errors.New(fmt.Sprintf("file size is greater than %d bytes\n", s))
		}
		file, err := os.Open(fname)
		if err != nil {
			return err
		}
		defer file.Close()
		r = file
	}

	return p.Print(r, w)
}
