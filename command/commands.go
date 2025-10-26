package command

import (
	"github.com/dihedron/archetype/command/initialise"
	"github.com/dihedron/archetype/command/show"
	"github.com/dihedron/archetype/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Initialise runs the Initialise command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Init initialise.Initialise `command:"initialise" alias:"init" alias:"i" description:"Initialise the project"`
	// Show runs the Show command which displays the settings needed for the specific project.
	Show show.Show `command:"show" alias:"s" description:"Show the necessary settings"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
