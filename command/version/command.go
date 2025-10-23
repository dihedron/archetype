package version

import (
	"log/slog"
	"os"

	"github.com/dihedron/repo/metadata"
)

// Version is the command that prints information about the application
// or plugin to the console; it support both compact and verbose mode.
type Version struct {
	Verbose bool `short:"v" long:"verbose" description:"Whether to print verbose information." optional:"yes"`
}

// Execute is the real implementation of the Version command.
func (cmd *Version) Execute(args []string) error {
	slog.Debug("running version command")
	if cmd.Verbose {
		metadata.PrintFull(os.Stdout)
	} else {
		metadata.Print(os.Stdout)
	}
	slog.Debug("command done")
	return nil
}
