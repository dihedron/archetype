package command

import (
	"github.com/dihedron/archetype/command/apply"
	"github.com/dihedron/archetype/command/describe"
	"github.com/dihedron/archetype/command/escape"
	"github.com/dihedron/archetype/command/unescape"
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
	// Escape runs the Escape command which escapes all Golang-template directives in the files in the repository.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Escape escape.Escape `command:"escape" alias:"esc" alias:"e" description:"Escape all Golang-template directives in the files in the repository"`
	// Escape runs the Escape command which escapes all Golang-template directives in the files in the repository.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Unescape unescape.Unescape `command:"unescape" alias:"unesc" alias:"u" description:"Unescape all Golang-template directives in the files in the repository"`
	// Version prints the application version information and exits.
	//lint:ignore SA5008 go-flags uses multiple tags to define aliases and choices
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit"`
}
