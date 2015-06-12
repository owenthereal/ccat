package main

import (
	"bytes"
	"testing"
)

func TestCPrint(t *testing.T) {
	r := bytes.NewBufferString("hello")
	var w bytes.Buffer

	err := CPrint(r, &w, LightColorPalettes)
	if err != nil {
		t.Errorf("error should be nil, but it's %s", err)
	}

	s := w.String()
	if s != "\033[34mhello\033[39;49;00m" {
		t.Errorf("output is wrong: %s", s)
	}
}
