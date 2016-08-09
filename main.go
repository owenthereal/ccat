package main

import (
	"fmt"
	"log"

	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

const (
	readFromStdin = "-"
)

type ccatCmd struct {
	BG          string
	Color       string
	ColorCodes  mapValue
	HTML        bool
	ShowPalette bool
	ShowVersion bool
}

func (c *ccatCmd) Run(cmd *cobra.Command, args []string) {
	stdout := colorable.NewColorableStdout()

	if c.ShowVersion {
		displayVersion(stdout)
		return
	}

	var colorPalettes ColorPalettes
	if c.BG == "dark" {
		colorPalettes = DarkColorPalettes
	} else {
		colorPalettes = LightColorPalettes
	}

	// override color codes
	for k, v := range c.ColorCodes {
		ok := colorPalettes.Set(k, v)
		if !ok {
			log.Fatal(fmt.Errorf("unknown color code: %s", k))
		}
	}

	if c.ShowPalette {
		fmt.Fprintf(stdout, `Applied color codes:

%s

Color code is in the format of:

  color       normal color
  *color*     bold color
  _color_     underlined color
  +color+     blinking color

Value of color can be %s
`, colorPalettes, colorCodes)
		return
	}

	var printer CCatPrinter
	if c.HTML {
		printer = HtmlPrinter{colorPalettes}
	} else if c.Color == "always" {
		printer = ColorPrinter{colorPalettes}
	} else if c.Color == "never" {
		printer = PlainTextPrinter{}
	} else {
		printer = AutoColorPrinter{colorPalettes}
	}

	// if there's no args, read from stdin
	if len(args) == 0 {
		args = []string{readFromStdin}
	}

	for _, arg := range args {
		err := CCat(arg, printer, stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	ccatCmd := &ccatCmd{
		ColorCodes: make(mapValue),
	}
	rootCmd := &cobra.Command{
		Use:  "ccat [OPTION]... [FILE]...",
		Long: "Colorize FILE(s), or standard input, to standard output.",
		Example: `$ ccat FILE1 FILE2 ...
  $ ccat --bg=dark FILE1 FILE2 ... # dark background
  $ ccat --html # output html
  $ ccat -G String="_darkblue_" -G Plaintext="darkred" FILE # set color codes
  $ ccat --palette # show palette
  $ ccat # read from standard input
  $ curl https://raw.githubusercontent.com/jingweno/ccat/master/main.go | ccat`,
		Run: ccatCmd.Run,
	}

	usageTempl := `{{ $cmd := . }}
Usage:
  {{.UseLine}}

Flags:
{{.LocalFlags.FlagUsages}}
Using color is auto both by default and with --color=auto. With --color=auto,
ccat emits color codes only when standard output is connected to a terminal.
Color codes can be changed with -G KEY=VALUE. List of color codes can
be found with --palette.

Examples:
  {{ .Example }}`
	rootCmd.SetUsageTemplate(usageTempl)

	rootCmd.PersistentFlags().StringVarP(&ccatCmd.BG, "bg", "", "light", `set to "light" or "dark" depending on the terminal's background`)
	rootCmd.PersistentFlags().StringVarP(&ccatCmd.Color, "color", "C", "auto", `colorize the output; value can be "never", "always" or "auto"`)
	rootCmd.PersistentFlags().VarP(&ccatCmd.ColorCodes, "color-code", "G", `set color codes`)
	rootCmd.PersistentFlags().BoolVarP(&ccatCmd.HTML, "html", "", false, `output html`)
	rootCmd.PersistentFlags().BoolVarP(&ccatCmd.ShowPalette, "palette", "", false, `show color palettes`)
	rootCmd.PersistentFlags().BoolVarP(&ccatCmd.ShowVersion, "version", "v", false, `show version`)

	rootCmd.Execute()
}
