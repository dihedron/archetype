package command

import (
	"github.com/dihedron/archetype/command/describe"
	"github.com/dihedron/archetype/command/generate"
	"github.com/dihedron/archetype/command/prepare"
	"github.com/dihedron/archetype/command/version"
)

// Commands is the main container for all the commands of the application.
type Commands struct {
	// Generate runs the Generate command
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Generate generate.Generate `command:"generate" alias:"init" alias:"apply" alias:"g" alias:"i" alias:"a" description:"Generate the project from the archetype"`
	// Describe runs the Describe command which displays the settings needed for the specific project.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Describe describe.Describe `command:"describe" alias:"descr" alias:"d" description:"Describe the necessary settings"`
	// Escape runs the Escape command which escapes all Golang-template directives in the files in the repository.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Escape prepare.Escape `command:"escape" alias:"esc" alias:"e" description:"Escape all Golang-template directives in the given files"`
	// Escape runs the Escape command which escapes all Golang-template directives in the files in the repository.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Unescape prepare.Unescape `command:"unescape" alias:"unesc" alias:"u" description:"Unescape all Golang-template directives in the given files"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
