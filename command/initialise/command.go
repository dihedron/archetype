package initialise

import (
	"bufio"
	"fmt"
	"log/slog"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/logging"
	"github.com/dihedron/archetype/repository"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type Initialise struct {
	base.Command
	// Settings is the path to the settings file to use for saturating the archetype variables.
	Settings Settings `short:"s" long:"settings" description:"The settings used to define the archetype and the specific parameters" required:"true"`
}

// Execute runs the Initialise command.
func (cmd *Initialise) Execute(args []string) error {
	slog.Info("executing Init command")

	var options []repository.Option

	fmt.Printf("%s\n", logging.ToYAML(cmd.Settings))

	if cmd.Settings.Version != 1 {
		slog.Error("unsupported settings version", "version", cmd.Settings.Version)
		return fmt.Errorf("unsupported settings version: %d", cmd.Settings.Version)
	}

	if cmd.Settings.Repository.URL == "" {
		slog.Error("repository URL not specified in settings")
		return fmt.Errorf("repository URL not specified in settings")
	}

	if cmd.Settings.Repository.Auth != nil {
		// TODO: validate auth settings
	}

	auth, err := cmd.AuthenticationOpts()
	if err != nil {
		slog.Error("error validating authentication options", "error", err)
		return err
	}
	if auth != nil {
		options = append(options, auth)
	}
	//options = append(options, repository.WithProxyFromEnv())

	// 1. create an in-memory clone of the remote archetypal repository
	repo := repository.New(
		cmd.Repository,
		options...,
	)
	repo.Clone()

	commit, err := repo.Commit(cmd.Tag)
	if err != nil {
		slog.Error("failed to get commit for input tag", "tag", cmd.Tag, "error", err)
		return err
	}

	/*
		// 4. load the parameters (TODO)
		for _, parameter := range cmd.Parameters.Values {
			fmt.Printf("%s %v\n", parameter.Name, parameter.Value)
		}
	*/

	// 5. loop over the files and perform some processing
	repo.ForEachFile(commit, visitFile)

	// 6. launch the script for post processing (TODO)

	return nil
}

func visitFile(file *object.File) error {
	fmt.Printf("%v  %9d  %s    %s\n", file.Mode, file.Size, file.Hash.String(), file.Name)

	reader, err := file.Blob.Reader()
	if err != nil {
		return err
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		slog.Error("error scanning file", "error", err)
		return err
	}

	return nil
}
