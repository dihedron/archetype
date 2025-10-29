package initialise

import (
	"bufio"
	"fmt"
	"log/slog"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/command/commons"
	"github.com/dihedron/archetype/logging"
	"github.com/dihedron/archetype/pointer"
	"github.com/dihedron/archetype/repository"
	"github.com/dihedron/archetype/settings"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type Initialise struct {
	base.Command
	// Settings is the path to the settings file to use for saturating the archetype variables.
	Settings settings.Settings `short:"s" long:"settings" description:"The settings used to define the archetype and the specific parameters" required:"true"`
}

const CurrentVersion = 1

// Execute runs the Initialise command.
func (cmd *Initialise) Execute(args []string) error {
	slog.Info("executing Initialise command")

	var options []repository.Option

	fmt.Printf("%s", logging.ToYAML(cmd.Settings))

	if cmd.Settings.Version != CurrentVersion {
		slog.Error("unsupported settings version", "version", cmd.Settings.Version, "expected", CurrentVersion)
		return fmt.Errorf("unsupported settings version: %d (expected %d)", cmd.Settings.Version, CurrentVersion)
	}

	if cmd.Settings.Repository.URL == "" {
		slog.Error("repository URL not specified in settings")
		return fmt.Errorf("repository URL not specified in settings")
	}

	// 1. extract authentication options
	if cmd.Settings.Repository.Auth != nil {
		// extract and validate auth settings
		auth, err := commons.AuthenticationOpts(cmd.Settings.Repository)
		if err != nil {
			slog.Error("error validating authentication options", "error", err)
			return err
		}
		if auth != nil {
			options = append(options, auth)
		}
	}
	//options = append(options, repository.WithProxyFromEnv())

	// 2. create an in-memory clone of the remote archetypal repository
	repo := repository.New(
		cmd.Settings.Repository.URL,
		options...,
	)
	if err := repo.Clone(); err != nil {
		slog.Error("failed to clone remote repository", "url", cmd.Settings.Repository.URL, "error", err)
		return err
	}

	// 3. checkout the specified tag (or 'latest' if none specified)
	if cmd.Settings.Repository.Tag == nil {
		slog.Info("no tag specified, using 'latest' as default")
		cmd.Settings.Repository.Tag = pointer.To("latest")
	}
	commit, err := repo.Commit(*cmd.Settings.Repository.Tag)
	if err != nil {
		slog.Error("failed to get commit for input tag", "tag", *cmd.Settings.Repository.Tag, "error", err)
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
