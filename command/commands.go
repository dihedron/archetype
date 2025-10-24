package command

import (
	"github.com/dihedron/archetype/command/bootstrap"
	"github.com/dihedron/archetype/command/show"
	"github.com/dihedron/archetype/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Bootstrap runs the Bootstrap command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Bootstrap bootstrap.Bootstrap `command:"bootstrap" alias:"boot" alias:"b" description:"Bootstrap the repository"`
	// Show runs the Show command which displays the settings needed for the specific project.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Show show.Show `command:"show" alias:"sh" alias:"s" description:"Show the necessary settings"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
