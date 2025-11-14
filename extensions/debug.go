package extensions

import (
	"encoding/json"
	"fmt"
	"os"
)

func toJSON(v any) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}

// DumpArgs is a template function that dumps the arguments passed to it to
// standard error.
func DumpArgs(args ...interface{}) (string, error) {
	result := ""
	if args != nil {
		for i, arg := range args {
			result += fmt.Sprintf("%d => '%s' (%T)\n", i, toJSON(arg), arg)
		}
		fmt.Fprintln(os.Stderr, result)
		return result, nil
	} else {
		return "<empty>", nil
	}
}
