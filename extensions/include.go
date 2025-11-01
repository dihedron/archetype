package extensions

import (
	"bufio"
	"errors"
	"strings"
	"text/template"
)

// Include is a template function that allows including a template from a file.
// It is roughly equivalent to the "template" keyword, but it also allows to
// specify an optional padding string that is prepended to each line of the
// included template.
// The function signature is:
//
//	include <template> [<pipeline>] [<padding>]
//
// where:
//   - <template> is the path to the template file to include.
//   - <pipeline> is an optional map[string]interface{} that is used as the
//     pipeline for the template execution.
//   - <padding> is an optional string that is prepended to each line of the
//     included template.
func Include(args ...interface{}) (string, error) {
	var (
		file    string
		padding string
		dynamic map[string]interface{}
	)
	if args == nil {
		return "", errors.New("include: at least the template path must be specified")
	}
	var pipelineFound bool
	for i, arg := range args {
		var ok bool

		if i == 0 {
			if file, ok = arg.(string); !ok {
				return "", errors.New("include: the first argument (template) must be of type string")
			}
		} else if i == 1 {
			if dynamic, ok = arg.(map[string]interface{}); !ok {
				if padding, ok = arg.(string); !ok {
					return "", errors.New("include: the second argument must either the pipeline or the padding")
				}
			} else {
				pipelineFound = true
			}
		} else if i == 2 {
			if !pipelineFound {
				return "", errors.New("include: the pipeline has not been provided")
			}
			if padding, ok = arg.(string); !ok {
				return "", errors.New("include: the third argument (padding) must be of type string")
			}
		}
	}

	// load the template
	t, err := template.ParseFiles(file)
	if err != nil {
		return "", err
	}

	var buffer strings.Builder
	if err = t.Execute(&buffer, dynamic); err != nil {
		return "", err
	}

	text := buffer.String()

	// apply padding only if necessary
	if padding != "" {
		var output strings.Builder
		scanner := bufio.NewScanner(strings.NewReader(text))
		for scanner.Scan() {
			output.WriteString(padding)
			output.WriteString(scanner.Text())
			output.WriteString("\n")
		}
		if scanner.Err() != nil {
			return "", err
		}
		return output.String(), nil
	}
	return text, nil
}
