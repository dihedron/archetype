package repository

import (
	"errors"
	"log/slog"

	"github.com/go-git/go-git/v6/plumbing"
)

// Branch retrieves a branch reference by its name.
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

// MainBranch retrieves the main branch reference, trying "main" first and
// falling back to "master" if "main" is not found.
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

// Branches retrieves all branch references.
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

// ForEachBranch iterates over all branches in the repository and applies the
// provided ReferenceVisitor function to each branch reference.
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
