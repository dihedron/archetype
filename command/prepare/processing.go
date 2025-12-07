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

// ProcessFile processes the given file; it returns a boolean indeicating
// whether any changes were made, and an error if any.
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
