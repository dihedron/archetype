package describe

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/dihedron/archetype/printf"
	"github.com/dihedron/archetype/repository"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// FileVisitor returns a function that processes files in a directory using the provided context for template rendering.
// The returned function is of type repository.FileVisitor, which is a callback that is invoked for each file in the repository.
// It skips files in the .archetype directory, and for all other files, it reads their content, parses them as text/template templates,
// executes them with the provided context, and writes the output to the corresponding path in the destination directory.
// It also adds the Sprig and custom template functions to the template.
func FileVisitor(excludePatterns []string, includePatterns []string) repository.FileVisitor {

	re := regexp.MustCompile(`(?s){{.*?}}`)

	includes := make([]*regexp.Regexp, 0)
	excludes := make([]*regexp.Regexp, 0)

	if len(includePatterns) > 0 {
		slog.Warn("include patterns are provided and will take precedence")
		for _, i := range includePatterns {
			slog.Info("including files matching pattern", "pattern", i)
			if re, err := regexp.Compile(i); err != nil {
				slog.Error("error compiling include pattern", "pattern", i, "error", err)
			} else {
				includes = append(includes, re)
			}
		}
	} else if len(excludePatterns) > 0 {
		for _, e := range excludePatterns {
			slog.Info("excluding files matching pattern", "pattern", e)
			if re, err := regexp.Compile(e); err != nil {
				slog.Error("error compiling exclude pattern", "pattern", e, "error", err)
			} else {
				excludes = append(excludes, re)
			}
		}
	}

	return func(file *object.File) error {
		// 1. skip files in the archive metadata directory
		if strings.HasPrefix(file.Name, ".archetype") {
			slog.Info("skipping archetype files", "file", file.Name)
			return nil
		}

		// fmt.Printf("exlcude: %d, include: %d\n", len(excludePatterns), len(includePatterns))

		//fmt.Printf("checking file %s\n", file.Name)
		if len(includes) > 0 {
			matched := false
			for _, re := range includes {
				if re.MatchString(file.Name) {
					matched = true
					break
				}
			}
			if !matched {
				slog.Info("skipping file not matching include patterns", "file", file.Name)
				//fmt.Printf("skipping file %s (no include pattern matches)\n", file.Name)
				return nil
			}
		} else if len(excludes) > 0 {
			for _, re := range excludes {
				if re.MatchString(file.Name) {
					slog.Info("skipping file matching exclude pattern", "file", file.Name)
					fmt.Printf("skipping file %s (exclude pattern matches)\n", file.Name)
					return nil
				}
			}
		}

		fmt.Printf("================================ %s ================================\n", printf.Green(file.Name))

		// 1. extract file contents into string
		text, err := file.Contents()
		if err != nil {
			fmt.Printf("%s getting file contents: %v\n", printf.Red("ERROR"), err)
			slog.Error("error getting file contents", "file", file.Name, "error", err)
			return err
		}

		// 2. Find all matches and their indexes.
		// FindAllStringIndex returns a slice of [start, end] pairs.
		// The -1 argument means "find all matches".
		matches := re.FindAllStringIndex(text, -1)

		// 3. Check if any matches were found.
		if len(matches) == 0 {
			slog.Debug("no template actions found in file", "file", file.Name)
			fmt.Println(text)
			return nil
		}

		slog.Debug("found matches in file", "file", file.Name, "matches", len(matches))

		// 4. Iterate through the string, printing in color.
		// We'll use lastIdx to keep track of where the last match ended.
		lastIdx := 0

		for _, match := range matches {
			// match[0] is the start index of the match
			// match[1] is the end index of the match

			// Print the part of the string *before* the match
			// This is text[lastIdx (end of last match) : start of current match]
			fmt.Print(text[lastIdx:match[0]])

			// Print the matched text in the highlight color
			// text[start of current match : end of current match]
			fmt.Print(printf.Magenta(text[match[0]:match[1]]))

			// Update our position to the end of the current match
			lastIdx = match[1]
		}

		// 5. Print any remaining text *after* the last match
		// This is the text from the end of the last match to the end of the string.
		fmt.Print(text[lastIdx:])

		// Add a final newline for clean terminal output
		//fmt.Println()

		fmt.Printf("-------------------------------- %s --------------------------------\n\n", printf.Green(file.Name))

		// --- Configuration ---
		// The regex pattern to search for.
		// We use `(?s)` to make the `.` character match newlines,
		// which is crucial for multi-line template actions.

		/*

			// process the filename as a template; the name of the file may be itself a template
			// and needs being renamed according to the values in the context; for instance, a file
			// named {{.ProjectName}}-config.yml should be rendered as myapp-config.yml if the
			// ProjectName in the context is "myapp"
			filename, err := template.New("filename").Parse(file.Name)
			if err != nil {
				slog.Error("cannot parse filename template", "template", file.Name, "error", err)
				return err
			}
			var buffer bytes.Buffer
			if err := filename.Execute(&buffer, context); err != nil {
				slog.Error("cannot execute filename template", "template", file.Name, "error", err)
				return err
			}
			output := path.Join(path.Clean(directory), buffer.String())
			fmt.Printf("%v  %9d  %s => ", file.Mode, file.Size, file.Name)
			//fmt.Printf("processing file %s (mode: %v, size: %d, hash: %s) as %s...\n", file.Name, file.Mode, file.Size, file.Hash.String(), output)
			//fmt.Printf("%s (mode: %v, size: %d): ", file.Name, file.Mode, file.Size)
			slog.Info("visiting file", "file", file.Name, "output", output, "mode", file.Mode, "size", file.Size, "output", output)

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
			buffer.Reset()
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
		*/
		return nil
	}
}
