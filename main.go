package main

import (
	"fmt"
	"log"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/mattn/go-colorable"
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/spf13/cobra"
)

const (
	readFromStdin = "-"
)

type ccatCmd struct {
	BG          string
	Color       string
	ColorCodes  mapValue
	ShowPalette bool
}

func (c *ccatCmd) Run(cmd *cobra.Command, args []string) {
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
		fmt.Println(colorPalettes)
		return
	}

	var printer CCatPrinter
	if c.Color == "always" {
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

	stdout := colorable.NewColorableStdout()
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
  $ ccat --color-code String="_darkblue_" --color-code Plaintext="darkred" FILE # set color codes
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
Color codes can be changed with --color-code KEY=VALUE.

Examples:
  {{ .Example }}`
	rootCmd.SetUsageTemplate(usageTempl)

	rootCmd.PersistentFlags().StringVarP(&ccatCmd.BG, "bg", "", "light", `set to "light" or "dark" depending on the terminal's background`)
	rootCmd.PersistentFlags().StringVarP(&ccatCmd.Color, "color", "C", "auto", `colorize the output; value can be "never", "always" or "auto"`)
	rootCmd.PersistentFlags().VarP(&ccatCmd.ColorCodes, "color-code", "G", `set color codes`)
	rootCmd.PersistentFlags().BoolVarP(&ccatCmd.ShowPalette, "palette", "", false, `show color palettes`)

	rootCmd.Execute()
}
