package main

import (
	"log/slog"
	"os"

	"github.com/dihedron/archetype/command"
	"github.com/jessevdk/go-flags"
)

func main() {

	defer cleanup()

	// if len(os.Args) == 2 && (os.Args[1] == "version" || os.Args[1] == "--version") {
	// 	metadata.Print(os.Stdout)
	// 	os.Exit(0)
	// } else if len(os.Args) == 3 && os.Args[1] == "version" && (os.Args[2] == "--verbose" || os.Args[2] == "-v") {
	// 	metadata.PrintFull(os.Stdout)
	// 	os.Exit(0)
	// }

	//var options Options
	options := command.Commands{}
	if _, err := flags.Parse(&options); err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			slog.Debug("help requested")
			os.Exit(0)
		}
		slog.Error("error parsing command line", "error", err)
		os.Exit(1)
	}

	/*
		auth, err := options.getAuthenticationOpts()
		if err != nil {
			slog.Error("error validating authentication options", "error", err)
			os.Exit(1)
		}

		repo := repository.New(
			options.Repository,
			auth,
		)
		repo.Clone()

		var reference *plumbing.Reference
		if options.Tag == "HEAD" {
			var err error
			// ... retrieving the branch being pointed by HEAD
			reference, err = repo.Head()
			if err != nil {
				slog.Error("failed to get HEAD", "error", err)
				os.Exit(1)
			}
			fmt.Println("HEAD points to:", reference.Name())
		} else {
			var err error
			reference, err = repo.Tag(options.Tag)
			if err != nil {
				slog.Error("failed to get tag", "error", err)
				os.Exit(1)
			}
		}

		repo.ForEachFile(reference, VisitFile)
	*/
}
