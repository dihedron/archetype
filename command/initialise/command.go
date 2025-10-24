package initialise

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"

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

	repo := repository.New(
		cmd.Repository,
		options...,
	)
	repo.Clone()

	longCommit := regexp.MustCompile(`(?m)^[0-9a-fA-F]{40}$`)
	shortCommit := regexp.MustCompile(`(?m)^[0-9a-fA-F]{7}$`)
	//var reference *plumbing.Reference
	var commit *object.Commit
	if cmd.Tag == "latest" || cmd.Tag == "HEAD" {
		// ... retrieving the branch being pointed by HEAD
		slog.Debug("retrieving reference to latest (HEAD)")
		reference, err := repo.Head()
		if err != nil {
			slog.Error("failed to get latest", "error", err)
			return err
		}
		fmt.Println("latest points to:", reference.Name())
		commit, err = repo.CommitFromReference(reference)
		if err != nil {
			slog.Error("failed to get commit for reference", "reference", reference.Name(), "error", err)
			return err
		}
		slog.Debug("retrieved commit for reference", "reference", reference.Name(), "hash", commit.Hash.String())
	} else if longCommit.MatchString(cmd.Tag) || shortCommit.MatchString(cmd.Tag) {
		slog.Debug("retrieving commit for specific hash", "hash", cmd.Tag)
		var err error
		commit, err = repo.Commit(cmd.Tag)
		if err != nil {
			slog.Error("failed to get commit", "error", err)
			return err
		}
		slog.Debug("retrieved commit for hash", "hash", cmd.Tag, "commit hash", commit.Hash.String())
	} else {
		var err error
		slog.Debug("retrieving reference to tag", "tag", cmd.Tag)
		reference, err := repo.Tag(cmd.Tag)
		if err != nil {
			slog.Error("failed to get tag", "error", err)
			os.Exit(1)
		}
		commit, err = repo.CommitFromReference(reference)
		if err != nil {
			slog.Error("failed to get commit for tag reference", "tag", cmd.Tag, "reference", reference.Name(), "error", err)
			return err
		}
		slog.Debug("retrieved commit for tag reference", "tag", cmd.Tag, "reference", reference.Name(), "hash", commit.Hash.String())
	}

	repo.ForEachFile(commit, visitFile)

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
