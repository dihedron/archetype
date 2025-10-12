package repository

import (
	"errors"
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
	h := plumbing.NewHash(hash)

	commit, err := r.repository.CommitObject(h)
	if err != nil {
		slog.Error("failed to get commit", "hash", hash, "error", err)
		return nil, err
	}
	return commit, nil
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
