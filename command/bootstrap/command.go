package bootstrap

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/dihedron/repo/command/base"
	"github.com/dihedron/repo/repository"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type Bootstrap struct {
	base.Command
}

func (cmd *Bootstrap) Execute(args []string) error {
	slog.Info("executing Bootstrap command")

	var options []repository.Option

	auth, err := cmd.AuthenticationOpts()
	if err != nil {
		slog.Error("error validating authentication options", "error", err)
		return err
	}
	if auth != nil {
		options = append(options, auth)
	}

	repo := repository.New(
		cmd.Repository,
		options...,
	)
	repo.Clone()

	var reference *plumbing.Reference
	if cmd.Tag == "HEAD" {
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
		reference, err = repo.Tag(cmd.Tag)
		if err != nil {
			slog.Error("failed to get tag", "error", err)
			os.Exit(1)
		}
	}

	repo.ForEachFile(reference, visitFile)

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
