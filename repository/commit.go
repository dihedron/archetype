package repository

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

// Commit retrieves a commit object by its hash.
func (r *Repository) Commit(hash string) (*object.Commit, error) {
	if r == nil || r.repository == nil {
		slog.Error("repository not initialized")
		return nil, errors.New("repository not initialized")
	}
	var h plumbing.Hash
	if len(hash) == 7 {
		slog.Debug("using short hash", "hash", hash)
		ref, err := r.repository.ResolveRevision(plumbing.Revision(hash))
		if err != nil {
			return nil, fmt.Errorf("error resolving revision '%s': %w", hash, err)
		}
		// convert the *plumbing.Reference to a plumbing.Hash
		h = *ref
	} else {
		slog.Debug("using long hash", "hash", hash)
		h = plumbing.NewHash(hash)
	}
	commit, err := r.repository.CommitObject(h)
	if err != nil {
		slog.Error("failed to get commit", "hash", hash, "error", err)
		return nil, err
	}
	return commit, nil
}

// CommitFromReference returns the Commit object that is pointed
// to by a given reference.
func (r *Repository) CommitFromReference(reference *plumbing.Reference) (*object.Commit, error) {
	// ... retrieving the commit object
	commit, err := r.repository.CommitObject(reference.Hash())
	if err != nil {
		slog.Error("error getting commit object for reference", "reference", reference.Name(), "error", err)
		return nil, err
	}
	return commit, err
}

// Head retrieves the current HEAD reference of the repository.
func (r *Repository) Head() (*plumbing.Reference, error) {
	if r == nil || r.repository == nil {
		slog.Error("repository not initialized")
		return nil, errors.New("repository not initialized")
	}
	reference, err := r.repository.Head()
	if err != nil {
		slog.Error("failed to get HEAD", "error", err)
		return nil, err
	}
	return reference, nil
}

// Commits retrieves all commit objects.
func (r *Repository) Commits() ([]*object.Commit, error) {
	var commits []*object.Commit
	err := r.ForEachCommit(func(commit *object.Commit) error {
		commits = append(commits, commit)
		return nil
	})
	if err != nil {
		slog.Error("failed to iterate commits", "error", err)
		return nil, err
	}
	return commits, nil
}

// ForEachCommit iterates over all commits in the repository and applies the
// provided CommitVisitor function to each commit.
func (r *Repository) ForEachCommit(visitor CommitVisitor) error {
	if r == nil || r.repository == nil {
		slog.Error("repository not initialized")
		return errors.New("repository not initialized")
	}
	commits, err := r.repository.CommitObjects()
	if err != nil {
		slog.Error("failed to get commits", "error", err)
		return err
	}
	err = commits.ForEach(visitor)
	if err != nil {
		slog.Error("failed to iterate branches", "error", err)
		return err
	}
	return nil
}
