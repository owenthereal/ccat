package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type HtmlCodes map[string]string

func (c HtmlCodes) String() string {
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
		s = append(s, Htmlize(ss, ss))
	}

	return strings.Join(s, ", ")
}

var htmlCodes = HtmlCodes{
	"":          "",
	"reset":     `</span>`,
	"bold":      `<span class="bold">`,
	"faint":     `<span class="faint">`,
	"standout":  `<span class="standout">`,
	"underline": `<span class="underline">`,
	"blink":     `<span class="blink">`,
	"overline":  `<span class="overline">`,
}

func init() {
	darkHtmls := []string{
		"black",
		"darkred",
		"darkgreen",
		"brown",
		"darkblue",
		"purple",
		"teal",
		"lightgray",
	}

	lightHtmls := []string{
		"darkgray",
		"red",
		"green",
		"yellow",
		"blue",
		"fuchsia",
		"turquoise",
		"white",
	}

	for i, x := 0, 30; i < len(darkHtmls); i, x = i+1, x+1 {
		htmlCodes[darkHtmls[i]] = fmt.Sprintf(`<span class="%s">`, darkHtmls[i])
		htmlCodes[lightHtmls[i]] = fmt.Sprintf(`<span class="%s">`, lightHtmls[i])
	}

	htmlCodes["darkteal"] = htmlCodes["turquoise"]
	htmlCodes["darkyellow"] = htmlCodes["brown"]
	htmlCodes["fuscia"] = htmlCodes["fuchsia"]
	htmlCodes["white"] = htmlCodes["bold"]
}

func Htmlize(attr, text string) string {
	if attr == "" {
		return text
	}

	result := new(bytes.Buffer)

	if strings.HasPrefix(attr, "+") && strings.HasSuffix(attr, "+") {
		result.WriteString(htmlCodes["blink"])
		attr = strings.TrimPrefix(attr, "+")
		attr = strings.TrimSuffix(attr, "+")
	}

	if strings.HasPrefix(attr, "*") && strings.HasSuffix(attr, "*") {
		result.WriteString(htmlCodes["bold"])
		attr = strings.TrimPrefix(attr, "*")
		attr = strings.TrimSuffix(attr, "*")
	}

	if strings.HasPrefix(attr, "_") && strings.HasSuffix(attr, "_") {
		result.WriteString(htmlCodes["underline"])
		attr = strings.TrimPrefix(attr, "_")
		attr = strings.TrimSuffix(attr, "_")
	}

	result.WriteString(htmlCodes[attr])
	result.WriteString(text)
	result.WriteString(htmlCodes["reset"])

	return result.String()
}
