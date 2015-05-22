package main

import (
	"bytes"
	"fmt"
	"strings"
)

const esc = "\033["

var ColorCodes = map[string]string{
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
		ColorCodes[darkColors[i]] = esc + fmt.Sprintf("%dm", x)
		ColorCodes[lightColors[i]] = esc + fmt.Sprintf("%d;01m", x)
	}

	ColorCodes["darkteal"] = ColorCodes["turquoise"]
	ColorCodes["darkyellow"] = ColorCodes["brown"]
	ColorCodes["fuscia"] = ColorCodes["fuchsia"]
	ColorCodes["white"] = ColorCodes["bold"]
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
		result.WriteString(ColorCodes["blink"])
		attr = strings.TrimPrefix(attr, "+")
		attr = strings.TrimSuffix(attr, "+")
	}

	if strings.HasPrefix(attr, "*") && strings.HasSuffix(attr, "*") {
		result.WriteString(ColorCodes["bold"])
		attr = strings.TrimPrefix(attr, "*")
		attr = strings.TrimSuffix(attr, "*")
	}

	if strings.HasPrefix(attr, "_") && strings.HasSuffix(attr, "_") {
		result.WriteString(ColorCodes["underline"])
		attr = strings.TrimPrefix(attr, "_")
		attr = strings.TrimSuffix(attr, "_")
	}

	result.WriteString(ColorCodes[attr])
	result.WriteString(text)
	result.WriteString(ColorCodes["reset"])

	return result.String()
}
