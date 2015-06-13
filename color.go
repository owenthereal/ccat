package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

const esc = "\033["

type ColorCodes map[string]string

func (c ColorCodes) String() string {
	var cc []string
	for k, _ := range c {
		if k == "" {
			continue
		}

		cc = append(cc, k)
	}
	sort.Strings(cc)

	var s []string
	for _, ss := range cc {
		s = append(s, Colorize(ss, ss))
	}

	return strings.Join(s, ", ")
}

var colorCodes = ColorCodes{
	"":          "",
	"reset":     esc + "39;49;00m",
	"bold":      esc + "01m",
	"faint":     esc + "02m",
	"standout":  esc + "03m",
	"underline": esc + "04m",
	"blink":     esc + "05m",
	"overline":  esc + "06m",
}

func init() {
	darkColors := []string{
		"black",
		"darkred",
		"darkgreen",
		"brown",
		"darkblue",
		"purple",
		"teal",
		"lightgray",
	}

	lightColors := []string{
		"darkgray",
		"red",
		"green",
		"yellow",
		"blue",
		"fuchsia",
		"turquoise",
		"white",
	}

	for i, x := 0, 30; i < len(darkColors); i, x = i+1, x+1 {
		colorCodes[darkColors[i]] = esc + fmt.Sprintf("%dm", x)
		colorCodes[lightColors[i]] = esc + fmt.Sprintf("%d;01m", x)
	}

	colorCodes["darkteal"] = colorCodes["turquoise"]
	colorCodes["darkyellow"] = colorCodes["brown"]
	colorCodes["fuscia"] = colorCodes["fuchsia"]
	colorCodes["white"] = colorCodes["bold"]
}

/*
	Format ``text`` with a color and/or some attributes::

		color       normal color
		*color*     bold color
		_color_     underlined color
		+color+     blinking color
*/
func Colorize(attr, text string) string {
	if attr == "" {
		return text
	}

	result := new(bytes.Buffer)

	if strings.HasPrefix(attr, "+") && strings.HasSuffix(attr, "+") {
		result.WriteString(colorCodes["blink"])
		attr = strings.TrimPrefix(attr, "+")
		attr = strings.TrimSuffix(attr, "+")
	}

	if strings.HasPrefix(attr, "*") && strings.HasSuffix(attr, "*") {
		result.WriteString(colorCodes["bold"])
		attr = strings.TrimPrefix(attr, "*")
		attr = strings.TrimSuffix(attr, "*")
	}

	if strings.HasPrefix(attr, "_") && strings.HasSuffix(attr, "_") {
		result.WriteString(colorCodes["underline"])
		attr = strings.TrimPrefix(attr, "_")
		attr = strings.TrimSuffix(attr, "_")
	}

	result.WriteString(colorCodes[attr])
	result.WriteString(text)
	result.WriteString(colorCodes["reset"])

	return result.String()
}
