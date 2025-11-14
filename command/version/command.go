package version

import (
	"log/slog"
	"os"

	"github.com/dihedron/archetype/metadata"
)

// Version is the command to print the version of the application.
type Version struct {
	Verbose bool `short:"v" long:"verbose" description:"Whether to print verbose information." optional:"yes"`
}

// Execute is the main entry point for the version command.
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
