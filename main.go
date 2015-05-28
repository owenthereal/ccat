package main

import (
	"log"

	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/mattn/go-colorable"
	"github.com/jingweno/ccat/Godeps/_workspace/src/github.com/spf13/cobra"
)

const (
	readFromStdin = "-"
)

type ccatCmd struct {
	BGFlag    string
	ColorFlag string
}

func (c *ccatCmd) Run(cmd *cobra.Command, args []string) {
	var colorDefs ColorDefs
	if c.BGFlag == "dark" {
		colorDefs = DarkColorDefs
	} else {
		colorDefs = LightColorDefs
	}

	var printer CCatPrinter
	if c.ColorFlag == "always" {
		printer = ColorPrinter{colorDefs}
	} else if c.ColorFlag == "never" {
		printer = PlainTextPrinter{}
	} else {
		printer = AutoColorPrinter{colorDefs}
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
	ccatCmd := &ccatCmd{}
	rootCmd := &cobra.Command{
		Use:  "ccat [OPTION]... [FILE]...",
		Long: "Colorize FILE(s), or standard input, to standard output.",
		Example: `$ ccat FILE1 FILE2 ...
  $ ccat --bg=dark FILE1 FILE2 ... # dark background
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

Examples:
  {{ .Example }}`
	rootCmd.SetUsageTemplate(usageTempl)

	rootCmd.PersistentFlags().StringVarP(&ccatCmd.BGFlag, "bg", "", "light", `Set to "light" or "dark" depending on the terminal's background`)
	rootCmd.PersistentFlags().StringVarP(&ccatCmd.ColorFlag, "color", "C", "auto", `colorize the output; value can be "never", "always" or "auto"`)

	rootCmd.Execute()
}
