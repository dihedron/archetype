package main

import (
	"log/slog"
	"os"

	"github.com/dihedron/archetype/command"
	"github.com/jessevdk/go-flags"
)

func main() {

	defer cleanup()

	options := command.Commands{}
	if _, err := flags.Parse(&options); err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			slog.Debug("help requested")
			os.Exit(0)
		}
		slog.Error("error parsing command line", "error", err)
		os.Exit(1)
	}
}
