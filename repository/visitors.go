package repository

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// ReferenceVisitor is a function type that defines the signature for visiting
// Git references. It takes a pointer to a plumbing.Reference and returns an error.
type ReferenceVisitor func(reference *plumbing.Reference) error

// DefaultReferenceVisitor is a sample implementation of ReferenceVisitor that
// logs the reference details.
func DefaultReferenceVisitor(reference *plumbing.Reference) error {
	slog.Info("visiting reference", "type", reference.Type().String(), "name", reference.Name().String())
	return nil
}

// CommitVisitor is a function type that defines the signature for visiting
// Git commits. It takes a pointer to an object.Commit and returns an error.
type CommitVisitor func(commit *object.Commit) error

// DefaultCommitVisitor is a sample implementation of CommitVisitor that
// logs the commit details.
func DefaultCommitVisitor(commit *object.Commit) error {
	slog.Info("visiting commit", "type", commit.Type().String(), "name", commit.Hash.String(), "message", commit.Message, "author", commit.Author.Name, "email", commit.Author.Email, "date", commit.Author.When.String())
	return nil
}

// FileVisitor is a function type that defines the signature for visiting
// Git files. It takes a pointer to an object.File and returns an error.
type FileVisitor func(file *object.File) error

// DefaultFileVisitor is a sample implementation of FileVisitor that
// logs the commit details.
func DefaultFileVisitor(file *object.File) error {
	slog.Info("visiting file", "type", file.Type().String(), "name", file.Name, "size", file.Size, "hash", file.Hash.String())
	fmt.Printf("%v  %9d  %s    %s\n", file.Mode, file.Size, file.Hash.String(), file.Name)
	return nil
}
