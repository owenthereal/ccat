package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/mattn/go-isatty"
)

func CCat(fname string, colorDefs ColorDefs) error {
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
		r = file
	}

	var err error
	if isatty.IsTerminal(uintptr(syscall.Stdout)) {
		err = CPrint(r, stdout, colorDefs)
	} else {
		_, err = io.Copy(stdout, r)
	}

	return err
}
