package repository

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// ReferenceVisitor is the signature of a function that can be used to visit
// a Git reference.
type ReferenceVisitor func(reference *plumbing.Reference) error

// DefaultReferenceVisitor is a sample implementation of the ReferenceVisitor that
// simply logs the reference details.
func DefaultReferenceVisitor(reference *plumbing.Reference) error {
	slog.Info("visiting reference", "type", reference.Type().String(), "name", reference.Name().String())
	return nil
}

// CommitVisitor is the signature of a function that can be used to visit
// a Git commit.
type CommitVisitor func(commit *object.Commit) error

// DefaultCommitVisitor is a sample implementation of the CommitVisitor that
// simply logs the commit details.
func DefaultCommitVisitor(commit *object.Commit) error {
	slog.Info("visiting commit", "type", commit.Type().String(), "name", commit.Hash.String(), "message", commit.Message, "author", commit.Author.Name, "email", commit.Author.Email, "date", commit.Author.When.String())
	return nil
}

// FileVisitor is the signature of a function that can be used to visit
// a Git file.
type FileVisitor func(file *object.File) error

// DefaultFileVisitor is a sample implementation of the FileVisitor that
// simply logs the file details.
func DefaultFileVisitor(file *object.File) error {
	slog.Info("visiting file", "type", file.Type().String(), "name", file.Name, "size", file.Size, "hash", file.Hash.String())
	fmt.Printf("%v  %9d  %s    %s\n", file.Mode, file.Size, file.Hash.String(), file.Name)
	return nil
}
