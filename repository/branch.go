package repository

import (
	"errors"
	"log/slog"

	"github.com/go-git/go-git/v6/plumbing"
)

// Branch returns the reference to the given branch.
func (r *Repository) Branch(name string) (*plumbing.Reference, error) {
	if r == nil || r.repository == nil {
		slog.Error("repository not initialized")
		return nil, errors.New("repository not initialized")
	}
	reference, err := r.repository.Reference(plumbing.NewBranchReferenceName(name), true)
	if err != nil {
		slog.Error("failed to get branch", "name", name, "error", err)
		return nil, err
	}
	return reference, nil
}

// MainBranch returns the reference to the main branch of the repository.
// It first tries to find a branch named "main", and if it does not exist,
// it falls back to "master".
func (r *Repository) MainBranch() (*plumbing.Reference, error) {
	reference, err := r.Branch("main")
	if errors.Is(err, plumbing.ErrReferenceNotFound) {
		slog.Debug("branch 'main' not found, trying 'master'...")
		// try "master" as a fallback
		reference, err = r.Branch("master")
		if err != nil {
			slog.Error("failed to get reference for master branch", "error", err)
			return nil, err
		}
	} else if err != nil {
		slog.Error("failed to get reference for main branch", "error", err)
		return nil, err
	}
	slog.Debug("main branch found", "name", reference.Name().String())
	return reference, nil
}

// Branches returns all the branches in the repository.
func (r *Repository) Branches() ([]*plumbing.Reference, error) {
	var branches []*plumbing.Reference
	err := r.ForEachBranch(func(reference *plumbing.Reference) error {
		branches = append(branches, reference)
		return nil
	})
	if err != nil {
		slog.Error("failed to iterate branches", "error", err)
		return nil, err
	}
	return branches, nil
}

// ForEachBranch iterates over all the branches in the repository and calls the
// visitor function for each branch.
func (r *Repository) ForEachBranch(visitor ReferenceVisitor) error {
	if r == nil || r.repository == nil {
		slog.Error("repository not initialized")
		return errors.New("repository not initialized")
	}
	branches, err := r.repository.Branches()
	if err != nil {
		slog.Error("failed to get branches", "error", err)
		return err
	}
	err = branches.ForEach(visitor)
	if err != nil {
		slog.Error("failed to iterate branches", "error", err)
		return err
	}
	return nil
}
