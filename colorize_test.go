package main

import "testing"

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
