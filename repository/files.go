package repository

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6/plumbing/object"
)

// Files returns the list of files in the given commit.
func (r *Repository) Files(commit *object.Commit) ([]*object.File, error) {

	if commit == nil {
		return nil, fmt.Errorf("invalid commit")
	}

	// ... retrieve the tree from the commit
	tree, err := commit.Tree()
	if err != nil {
		slog.Error("error getting tree for commit", "commit", commit.Hash, "error", err)
		return nil, err
	}

	files := []*object.File{}
	// ... get the files iterator and print the file
	tree.Files().ForEach(func(f *object.File) error {
		files = append(files, f)
		return nil
	})
	return files, nil
}

// ForEachFile iterates over all the files in the given commit and calls the
// visitor function for each file.
func (r *Repository) ForEachFile(commit *object.Commit, visitor FileVisitor) error {
	files, err := r.Files(commit)
	if err != nil {
		slog.Error("error getting files for commit", "hash", commit.Hash.String(), "error", err)
		return err
	}
	for _, file := range files {
		visitor(file)
	}
	return nil
}
