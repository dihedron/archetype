package prepare

import (
	"errors"
	"log/slog"
)

// Unescape is the command to unescape all Golang-template directives in a file.
type Unescape struct {
	// Directory is the path to the directory to use to store the "escaped" files.
	Directory string `short:"d" long:"directory" description:"The directory where the output files are stored" required:"true" default:".archetype/escaped"`
}

// Execute is the main entry point for the unescape command.
// It unescapes all Golang-template directives in the given files.
// The resulting files are written to the output directory.
func (cmd *Unescape) Execute(args []string) error {
	slog.Info("executing Unescape command")

	var errs error
	for _, arg := range args {

		if err := UnescapeFile(arg, cmd.Directory); err != nil {
			slog.Error("error unescaping file", "file", arg, "error", err)
			errs = errors.Join(errs, err)
		}
	}

	return errs
}
