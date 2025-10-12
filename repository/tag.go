package repository

import (
	"errors"
	"log/slog"

	"github.com/go-git/go-git/v6/plumbing"
)

// Tag retrieves a tag reference by its name.
func (r *Repository) Tag(name string) (*plumbing.Reference, error) {
	if r == nil || r.repository == nil {
		slog.Error("repository not initialized")
		return nil, errors.New("repository not initialized")
	}
	reference, err := r.repository.Reference(plumbing.NewTagReferenceName(name), true)
	if err != nil {
		slog.Error("failed to get tag", "name", name, "error", err)
		return nil, err
	}
	return reference, nil
}

// Tags retrieves all tag references.
func (r *Repository) Tags() ([]*plumbing.Reference, error) {
	var tags []*plumbing.Reference
	err := r.ForEachTag(func(reference *plumbing.Reference) error {
		tags = append(tags, reference)
		return nil
	})
	if err != nil {
		slog.Error("failed to iterate tags", "error", err)
		return nil, err
	}
	return tags, nil
}

// ForEachTag iterates over all tag references applying the visitor function.
func (r *Repository) ForEachTag(visitor ReferenceVisitor) error {
	if r == nil || r.repository == nil {
		slog.Error("repository not initialized")
		return errors.New("repository not initialized")
	}
	tags, err := r.repository.Tags()
	if err != nil {
		slog.Error("failed to get tags", "error", err)
		return err
	}
	err = tags.ForEach(visitor)
	if err != nil {
		slog.Error("failed to iterate tags", "error", err)
		return err
	}
	return nil
}
