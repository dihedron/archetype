package prepare

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/printf"
)

var (
	RealBra = "{{"
	SafeBra = "{-{"
	RealKet = "}}"
	SafeKet = "}-}"
)

// EscapeFile escapes all Golang-template directives in the given file.
func EscapeFile(filename string, directory string) error {
	_, err := ProcessFile(filename, directory, RealBra, RealKet, SafeBra, SafeKet)
	return err
}

// UnescapeFile unescapes all Golang-template directives in the given file.
func UnescapeFile(filename string, directory string) error {
	_, err := ProcessFile(filename, directory, SafeBra, SafeKet, RealBra, RealKet)
	return err
}

const (
	DefaultDirectoryPermissions = 0755
	DefaultFilePermissions      = 0644
)

// ProcessFile processes the given file; it returns a boolean indicating
// whether any changes were made, and an error if any occurred.
func ProcessFile(filename, directory, oldBra, oldKet, newBra, newKet string) (bool, error) {

	modified := false

	// 1. check if the input file exists
	if _, err := os.Stat(filename); err != nil {
		slog.Error("error getting file info", "file", filename, "error", err)
		return false, err
	}

	// 2. create the name of the output file
	output := path.Join(path.Clean(directory), filename)
	slog.Debug("processing file", "file", filename, "output", output)

	// 3. create the output directory structure if it does not exist
	err := os.MkdirAll(filepath.Dir(output), 0770)
	if err != nil {
		slog.Error("error creating directory structure", "path", filepath.Dir(output), "error", err)
		return false, err
	}

	// 4. check if file is binary and skip processing if so
	if isText, err := base.IsTextFile(filename); err != nil {
		slog.Error("error getting file info", "file", filename, "error", err)
		return false, err
	} else if !isText {
		slog.Debug("file is binary, skipping processing and copying as is", "file", filename, "output", output)
		if err := CopyFile(filename, output); err != nil {
			slog.Error("error copying file", "file", filename, "output", output, "error", err)
			return false, err
		}
		return false, nil
	} else {

		var data []byte

		// 5. extract file contents into string
		tmp, err := os.ReadFile(filename)
		if err != nil {
			slog.Error("error getting file contents", "file", filename, "error", err)
			return false, err
		}

		text := string(tmp)

		// 6. find all matches and their indexes.
		re := regexp.MustCompile(`(?s)` + oldBra + `.*?` + oldKet)
		matches := re.FindAllStringIndex(text, -1)

		// 7. check if any matches were found.
		if len(matches) == 0 {
			slog.Debug("no template actions found in file", "file", filename)
		} else {
			slog.Debug("found matches in file", "file", filename, "matches", len(matches))

			// 8. iterate through the string, printing in color.
			lastIdx := 0
			var buffer bytes.Buffer
			for _, match := range matches {
				// match[0] is the start index of the match
				// match[1] is the end index of the match

				// print the part of the string *before* the match
				fmt.Fprint(&buffer, text[lastIdx:match[0]])

				// print the matched text in the highlight color
				//fmt.Fprint(&buffer, printf.Magenta(text[match[0]:match[1]]))
				fmt.Fprint(&buffer, strings.Replace(strings.Replace(text[match[0]:match[1]], oldBra, newBra, 1), oldKet, newKet, 1))

				// Update our position to the end of the current match
				lastIdx = match[1]
			}

			// 9. Print any remaining text *after* the last match
			// This is the text from the end of the last match to the end of the string.
			fmt.Fprint(&buffer, text[lastIdx:])

			modified = true
			data = buffer.Bytes()
			slog.Debug("processed file", "file", filename)
		}

		// 10. output the rendered content
		if err := os.WriteFile(output, data, os.FileMode(DefaultFilePermissions)); err != nil {
			slog.Error("error writing file", "file", output, "error", err)
			return false, fmt.Errorf("error writing file %s: %w", output, err)
		}

	}
	return modified, nil
}

type BracketsProcessorFunc func(bra, ket, text string) string

// ProcessBrackets processes all occurrences of text between brackets.
func ProcessBrackets(data []byte, bra, ket string, process BracketsProcessorFunc) []byte {

	text := string(data)

	// 1. find all matches and their indexes.
	re := regexp.MustCompile(`(?s)` + bra + `.*?` + ket)
	matches := re.FindAllStringIndex(text, -1)

	// 2. check if any matches were found.
	if len(matches) == 0 {
		slog.Debug("no template actions found in data")
		return data
	} else {
		slog.Debug("found matches in data", "matches", len(matches))

		// 3. iterate through the data, processing the matched text.
		lastIdx := 0
		var buffer bytes.Buffer
		for _, match := range matches {
			// match[0] is the start index of the match
			// match[1] is the end index of the match

			// print the part of the string *before* the match
			fmt.Fprint(&buffer, text[lastIdx:match[0]])

			// process the matched text
			fmt.Fprint(&buffer, process(bra, ket, text[match[0]+len(bra):match[1]-len(ket)]))

			// update our position to the end of the current match
			lastIdx = match[1]
		}

		// 4. print any remaining text *after* the last match
		// this is the text from the end of the last match to the end of the string.
		fmt.Fprint(&buffer, text[lastIdx:])

		data = buffer.Bytes()
		slog.Debug("processed data")
	}

	return data
}

func BracketsReplacer(oldBra, oldKet, newBra, newKet string) BracketsProcessorFunc {
	return func(_, _, text string) string {
		return newBra + text + newKet
	}
}

func SelectiveBracketsReplacer(oldBra, oldKet, newBra, newKet string, accept func(string) bool) BracketsProcessorFunc {
	return func(_, _, text string) string {
		if accept(text) {
			return newBra + text + newKet
		}
		return oldBra + text + oldKet
	}
}

func BracketsColoriser(oldBra, oldKet, newBra, newKet string) BracketsProcessorFunc {
	return func(_, _, text string) string {
		return fmt.Sprintf("%s%s%s", printf.Magenta(oldBra), printf.Magenta(text), printf.Magenta(oldKet))
	}
}

func SelectiveBracketsColoriser(oldBra, oldKet, newBra, newKet string, accept func(string) bool) BracketsProcessorFunc {
	return func(_, _, text string) string {
		if accept(text) {
			return fmt.Sprintf("%s%s%s", printf.Magenta(oldBra), printf.Magenta(text), printf.Magenta(oldKet))
		}
		return fmt.Sprintf("%s%s%s", oldBra, text, oldKet)
	}
}

func ReplaceAllBrackets(data []byte, oldBra, oldKet, newBra, newKet string) []byte {
	return ProcessBrackets(data, oldBra, oldKet, BracketsReplacer(oldBra, oldKet, newBra, newKet))
}

func ReplaceSelectedBrackets(data []byte, oldBra, oldKet, newBra, newKet string, accept func(string) bool) []byte {
	return ProcessBrackets(data, oldBra, oldKet, SelectiveBracketsReplacer(oldBra, oldKet, newBra, newKet, accept))
}

func ColoriseAllBrackets(data []byte, oldBra, oldKet, newBra, newKet string) []byte {
	return ProcessBrackets(data, oldBra, oldKet, BracketsColoriser(oldBra, oldKet, newBra, newKet))
}

func ColoriseSelectedBrackets(data []byte, oldBra, oldKet, newBra, newKet string, accept func(string) bool) []byte {
	return ProcessBrackets(data, oldBra, oldKet, SelectiveBracketsColoriser(oldBra, oldKet, newBra, newKet, accept))
}

// CopyFile copies a file as is from source to destination.
func CopyFile(source, destination string) error {
	// 1. open the source file for reading
	reader, err := os.Open(source)
	if err != nil {
		slog.Error("error opening source file", "file", source, "error", err)
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer reader.Close()

	// 2. create (or truncate) the destination file
	writer, err := os.Create(destination)
	if err != nil {
		slog.Error("error creating destination file", "file", destination, "error", err)
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer writer.Close()

	// 3. copy the bytes from source to destination using io.Copy
	count, err := io.Copy(writer, reader)
	if err != nil {
		slog.Error("error during copy operation", "source", source, "destination", destination, "error", err)
		return fmt.Errorf("error during copy operation: %w", err)
	}
	slog.Debug("copied file", "source", source, "destination", destination, "bytes", count)

	// 4. ensure data is flushed to disk
	err = writer.Sync()
	if err != nil {
		slog.Error("error syncing to disk", "file", destination, "error", err)
		return fmt.Errorf("error syncing to disk: %w", err)
	}
	slog.Debug("synced file", "file", destination)

	return nil
}
