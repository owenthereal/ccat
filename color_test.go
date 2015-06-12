package main

import (
	"testing"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/sourcegraph/syntaxhighlight"
)

func Test_ColorDefs_Set(t *testing.T) {
	defs := ColorDefs{
		syntaxhighlight.String: "blue",
	}

	ok := defs.Set("foo", "bar")
	if ok {
		t.Errorf("setting color code foo should not be ok")
	}

	ok = defs.Set("String", "baz")
	if !ok {
		t.Errorf("setting color code String should be ok")
	}

	if defs[syntaxhighlight.String] != "baz" {
		t.Errorf("color code of String should be baz")
	}
}

func TestColorize(t *testing.T) {
	cases := []struct {
		Color, Output string
	}{
		{
			Color:  "",
			Output: "hello",
		},

		{
			Color:  "blue",
			Output: "\033[34;01mhello\033[39;49;00m",
		},
		{
			Color:  "_blue_",
			Output: "\033[04m\033[34;01mhello\033[39;49;00m",
		},
		{
			Color:  "bold",
			Output: "\033[01mhello\033[39;49;00m",
		},
	}

	for _, tc := range cases {
		actual := Colorize(tc.Color, "hello")
		if actual != tc.Output {
			t.Errorf(
				"Color: %#v\n\nOutput: %#v\n\nExpected: %#v",
				tc.Color,
				actual,
				tc.Output)
		}
	}
}

func TestColorizeMultiByte(t *testing.T) {
	cases := []struct {
		Color, Output string
	}{
		// Japanese
		{
			Color:  "",
			Output: "こんにちは",
		},

		{
			Color:  "blue",
			Output: "\033[34;01mこんにちは\033[39;49;00m",
		},
		{
			Color:  "_blue_",
			Output: "\033[04m\033[34;01mこんにちは\033[39;49;00m",
		},
		{
			Color:  "bold",
			Output: "\033[01mこんにちは\033[39;49;00m",
		},
	}

	for _, tc := range cases {
		actual := Colorize(tc.Color, "こんにちは")
		if actual != tc.Output {
			t.Errorf(
				"Color: %#v\n\nOutput: %#v\n\nExpected: %#v",
				tc.Color,
				actual,
				tc.Output)
		}
	}
}
