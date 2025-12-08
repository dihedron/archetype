package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dihedron/archetype/command"
	"github.com/dihedron/archetype/command/prepare"
	"github.com/jessevdk/go-flags"
)

func TestColorise() {
	expected := `
{{"\"output\""}}
	A string constant.
{{"output"}}
	A raw string constant.
{{printf "%q" "output"}}
	A function call.
{{"output" | printf "%q"}}
	A function call whose final argument comes from the previous
	command.
{{printf "%q" (print "out" "put")}}
	A parenthesized argument.
{{"put" | printf "%s%s" "out" | printf "%q"}}
	A more elaborate call.
{{"output" | printf "%s" | printf "%q"}}
	A longer chain.
{{with "output"}}{{printf "%q" .}}{{end}}
	A with action using dot.
{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
	A with action that creates and uses a variable.
{{with $x := "output"}}{{printf "%q" $x}}{{end}}
	A with action that uses the variable in another action.
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
	The same, but pipelined.	
`
	got := string(
		prepare.ColoriseSelectedBrackets(
			[]byte(expected),
			prepare.RealBra,
			prepare.RealKet,
			prepare.SafeBra,
			prepare.SafeKet,
			func(s string) bool {
				return strings.Contains(s, "output")
			},
		),
	)
	fmt.Println(got)
}

// main is the entry point of the application.
func main() {

	// TestColorise()

	// return

	defer cleanup()

	options := command.Commands{}
	if _, err := flags.Parse(&options); err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			slog.Debug("help requested")
			os.Exit(0)
		}
		slog.Error("error parsing command line", "error", err)
		os.Exit(1)
	}
}
