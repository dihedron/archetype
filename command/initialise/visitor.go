package initialise

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/dihedron/archetype/extensions"
	"github.com/dihedron/archetype/repository"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v6/plumbing/object"
)

var (
	green   func(...any) string = color.New(color.FgGreen).SprintFunc()
	red     func(...any) string = color.New(color.FgRed).SprintFunc()
	yellow  func(...any) string = color.New(color.FgYellow).SprintFunc()
	blue    func(...any) string = color.New(color.FgBlue).SprintFunc()
	magenta func(...any) string = color.New(color.FgMagenta).SprintFunc()
)

// FileVisitor creates a FileVisitor function that processes files
// in the given directory using the provided context for template rendering.
func FileVisitor(directory string, context any) repository.FileVisitor {

	return func(file *object.File) error {

		if strings.HasPrefix(file.Name, ".archetype") {
			slog.Info("skipping archetype files", "file", file.Name)
			return nil
		}

		output := path.Join(path.Clean(directory), file.Name)
		fmt.Printf("%v  %9d  %s => ", file.Mode, file.Size, file.Name)
		//fmt.Printf("processing file %s (mode: %v, size: %d, hash: %s) as %s...\n", file.Name, file.Mode, file.Size, file.Hash.String(), output)
		//fmt.Printf("%s (mode: %v, size: %d): ", file.Name, file.Mode, file.Size)
		slog.Info("visiting file", "file", file.Name)

		reader, err := file.Blob.Reader()
		if err != nil {
			fmt.Printf("%s getting file reader: %v\n", red("ERROR"), err)
			slog.Error("error getting file reader", "file", file.Name, "error", err)
			return err
		}
		defer reader.Close()

		contents, err := file.Contents()
		if err != nil {
			fmt.Printf("%s getting file contents: %v\n", red("ERROR"), err)
			slog.Error("error getting file contents", "file", file.Name, "error", err)
			return err
		}

		// populate the functions map
		functions := template.FuncMap{}
		for k, v := range extensions.FuncMap() {
			functions[k] = v
		}
		for k, v := range sprig.FuncMap() {
			functions[k] = v
		}

		// parse the templates
		main := path.Base(file.Name)
		templates, err := template.New(main).Funcs(functions).Parse(contents)
		if err != nil {
			slog.Error("cannot parse template file", "file", file.Name, "error", err)
			fmt.Printf("%s parsing template: %v\n", red("ERROR"), err)
			return fmt.Errorf("error parsing template file %v: %w", file.Name, err)
		}

		// execute the template
		var buffer bytes.Buffer
		if err := templates.ExecuteTemplate(&buffer, main, context); err != nil {
			slog.Error("cannot apply data to template", "error", err)
			fmt.Printf("%s applying data to template: %v\n", red("ERROR"), err)
			return fmt.Errorf("error applying data to template: %w", err)
		}

		// output the rendered content
		if err = os.WriteFile(output, buffer.Bytes(), os.FileMode(file.Mode)); err != nil {
			slog.Error("error writing file", "file", file.Name, "error", err)
			fmt.Printf("%s writing file as %s: %v\n", red("ERROR"), output, err)
			return fmt.Errorf("error writing file %s: %w", file.Name, err)
		}
		fmt.Printf("%s (as %s)\n", green("SUCCESS"), output)
		//fmt.Printf("---- rendered content of %s ----\n%s\n---- end of rendered content of %s ----\n", file.Name, buffer.String(), file.Name)
		return nil
	}
}
