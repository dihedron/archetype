package command

import (
	"github.com/dihedron/archetype/command/bootstrap"
	"github.com/dihedron/archetype/command/describe"
	"github.com/dihedron/archetype/command/version"
)

// Commands is the main container for all the commands of the application.
type Commands struct {
	// Bootstrap runs the Bootstrap command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Bootstrap bootstrap.Bootstrap `command:"bootstrap" alias:"boot" alias:"b" description:"Bootstrap the directory based on the archetype"`
	// Describe runs the Describe command which displays the settings needed for the specific project.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Describe describe.Describe `command:"describe" alias:"descr" alias:"d" description:"Describe the necessary settings"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
