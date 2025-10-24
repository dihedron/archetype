package repository

import (
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6/plumbing/object"
)

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

// ForEachFile iterates over all files in the repository and applies the
// provided FileVisitor function to each file.
// func (r *Repository) ForEachFile(reference *plumbing.Reference, visitor FileVisitor) error {
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
	// // ... retrieving the commit object
	// commit, err := r.repository.CommitObject(reference.Hash())
	// if err != nil {
	// 	slog.Error("error getting commit object for reference", "reference", reference.Name(), "error", err)
	// 	return err
	// }
	// fmt.Println(commit)

	// // ... retrieve the tree from the commit
	// tree, err := commit.Tree()
	// if err != nil {
	// 	slog.Error("error getting tree for commit", "commit", commit.Hash, "error", err)
	// 	return err
	// }

	// // ... get the files iterator and print the file
	// return tree.Files().ForEach(visitor)
}
