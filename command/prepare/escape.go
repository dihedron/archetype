package prepare

import (
	"errors"
	"log/slog"
)

// Escape is the command to escape all Golang-template directives in a file.
type Escape struct {
	// Directory is the path to the directory to use to store the "escaped" files.
	Directory string `short:"d" long:"directory" description:"The directory where the output files are stored" required:"true" default:".archetype/escaped"`
}

// Execute is the main entry point for the escape command.
// It escapes all Golang-template directives in the given files.
// The resulting files are written to the output directory.
func (cmd *Escape) Execute(args []string) error {
	slog.Info("executing Escape command")

	var errs error
	for _, arg := range args {

		if err := EscapeFile(arg, cmd.Directory); err != nil {
			slog.Error("error escaping file", "file", arg, "error", err)
			errs = errors.Join(errs, err)
		}
	}

	return errs
}
