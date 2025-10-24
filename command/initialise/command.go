package initialise

import (
	"bufio"
	"fmt"
	"log/slog"

	"github.com/dihedron/archetype/command/base"
	"github.com/dihedron/archetype/repository"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type Init struct {
	base.Command
}

func (cmd *Init) Execute(args []string) error {
	slog.Info("executing Init command")

	var options []repository.Option

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

	// 4. load the parameters (TODO)

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
