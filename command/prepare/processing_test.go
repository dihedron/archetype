package prepare

import "testing"

func TestReplaceBrackets(t *testing.T) {
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
		ReplaceAllBrackets(
			ReplaceAllBrackets(
				[]byte(expected),
				RealBra,
				RealKet,
				SafeBra,
				SafeKet,
			),
			SafeBra,
			SafeKet,
			RealBra,
			RealKet,
		),
	)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestPrintBrackets(t *testing.T) {
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
		ColoriseAllBrackets(
			[]byte(expected),
			RealBra,
			RealKet,
			SafeBra,
			SafeKet,
		),
	)
	t.Logf("expected %q", expected)
	t.Logf("got %q", got)
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
