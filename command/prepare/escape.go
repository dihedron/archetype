package prepare

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/dihedron/archetype/command/base"
)

// Escape is the command to escape all Golang-template directives in a file.
type Escape struct {
	// Directory is the path to the directory to use to store the "escaped" files.
	Directory string `short:"d" long:"directory" description:"The directory where the output files are stored" required:"true" default:".archetype/escaped"`
}

// // Execute is the main entry point for the escape command.
// // It escapes all Golang-template directives in the given files.
// // The resulting files are written to the output directory.
// func (cmd *Escape) ExecuteOld(args []string) error {
// 	slog.Info("executing Escape command")

// 	var errs error
// 	for _, arg := range args {

// 		if err := EscapeFile(arg, cmd.Directory); err != nil {
// 			slog.Error("error escaping file", "file", arg, "error", err)
// 			errs = errors.Join(errs, err)
// 		}
// 	}

// 	return errs
// }

// Execute is the main entry point for the escape command.
// It escapes all Golang-template directives in the given files.
// The resulting files are written to the output directory.
func (cmd *Escape) Execute(args []string) error {
	slog.Info("executing Escape command")

	var errs error
	for _, filename := range args {

		// 1. check if the input file exists
		if _, err := os.Stat(filename); err != nil {
			slog.Error("error getting file info", "file", filename, "error", err)
			return err
		}

		// 2. create the name of the output file
		output := path.Join(path.Clean(cmd.Directory), filename)
		slog.Debug("processing file", "file", filename, "output", output)

		// 3. create the output directory structure if it does not exist
		err := os.MkdirAll(filepath.Dir(output), DefaultDirectoryPermissions)
		if err != nil {
			slog.Error("error creating directory structure", "path", filepath.Dir(output), "error", err)
			return err
		}

		// 4. read the file contents into memory
		data, err := os.ReadFile(filename)
		if err != nil {
			slog.Error("error getting file contents", "file", filename, "error", err)
			return err
		}

		// 5. process the file contents, if it's a text file
		if base.IsText(data) {
			slog.Debug("file is text", "file", filename)
			data = ReplaceSelectedBrackets(data, RealBra, RealKet, SafeBra, SafeKet, func(s string) bool {
				// TODO: implement the selection logic
				return true
			})
		} else {
			slog.Debug("file is binary, skipping", "file", filename)
		}

		// 6. write out the file contents
		if err := os.WriteFile(output, data, DefaultFilePermissions); err != nil {
			slog.Error("error writing file", "file", output, "error", err)
			return err
		}
	}

	return errs
}
