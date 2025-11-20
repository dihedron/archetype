package command

import (
	"github.com/dihedron/archetype/command/apply"
	"github.com/dihedron/archetype/command/describe"
	"github.com/dihedron/archetype/command/prepare"
	"github.com/dihedron/archetype/command/version"
)

// Commands is the main container for all the commands of the application.
type Commands struct {
	// Apply runs the Apply command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Apply apply.Apply `command:"apply" alias:"init" alias:"a" description:"Apply the archetype to the project"`
	// Describe runs the Describe command which displays the settings needed for the specific project.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Describe describe.Describe `command:"describe" alias:"descr" alias:"d" description:"Describe the necessary settings"`
	// Prepare runs the Prepare command which sets up the necessary files and directories.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Prepare prepare.Prepare `command:"prepare" alias:"prep" alias:"p" description:"Prepare the template files by escaping {{ and }}"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
