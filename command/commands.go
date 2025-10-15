package command

import (
	"github.com/dihedron/bootstrap/command/apply"
	"github.com/dihedron/bootstrap/command/show"
	"github.com/dihedron/bootstrap/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Apply runs the Apply command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Apply apply.Apply `command:"apply" alias:"ap" alias:"a" description:"Bootstrap the project."`
	// Show runs the Show command which displays the settings needed for the specific project.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Show show.Show `command:"show" alias:"sh" alias:"s" description:"Show the necessary settings"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
