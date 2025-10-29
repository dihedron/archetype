package commons

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/dihedron/archetype/repository"
	"github.com/dihedron/archetype/settings"
)

// AuthenticationOpts extracts authentication data from a repository.Settings
// and creates the repository.Option needed to configure authenticated requests
// against the remote repository.
func AuthenticationOpts(repo settings.Repository) (repository.Option, error) {
	if repo.Auth == nil {
		slog.Info("using anonymous authentication")
		return nil, nil
	}
	if repo.Auth.Token != nil {
		if strings.HasPrefix(repo.URL, "http") {
			slog.Info("using token for authentication")
			return repository.WithTokenAuth(*repo.Auth.Token), nil
		} else {
			slog.Error("token authentication is only supported for HTTP repositories")
			return nil, errors.New("token authentication is only supported for HTTP repositories")
		}
	} else if repo.Auth.Password != nil && repo.Auth.Username != nil {
		if strings.HasPrefix(repo.URL, "http") {
			slog.Info("using username and password for authentication")
			return repository.WithBasicAuth(*repo.Auth.Username, *repo.Auth.Password), nil
		} else {
			slog.Error("username and password authentication is only supported for HTTP repositories")
			return nil, errors.New("username and password authentication is only supported for HTTP repositories")
		}
	} else if repo.Auth.SSHKey != nil {
		if strings.HasPrefix(repo.URL, "ssh://") || strings.HasPrefix(repo.URL, "git@") {
			slog.Info("using SSH key for authentication")
			return repository.WithSSHKey(*repo.Auth.SSHKey, nil), nil
		} else {
			slog.Error("SSH key authentication is only supported for SSH repositories")
			return nil, errors.New("SSH key authentication is only supported for SSH repositories")
		}
	} else if repo.Auth.UseDefaultSSHKey {
		if strings.HasPrefix(repo.URL, "ssh://") || strings.HasPrefix(repo.URL, "git@") {
			slog.Info("using default SSH key for authentication")
			return repository.WithDefaultSSHKey(), nil
		} else {
			slog.Error("SSH key authentication is only supported for SSH repositories")
			return nil, errors.New("SSH key authentication is only supported for SSH repositories")
		}
	} else if repo.Auth.UseSSHAgent {
		if strings.HasPrefix(repo.URL, "ssh://") || strings.HasPrefix(repo.URL, "git@") {
			slog.Info("using default SSH key for authentication")
			return repository.WithDefaultSSHKey(), nil
		} else {
			slog.Error("SSH agent authentication is only supported for SSH repositories")
			return nil, errors.New("SSH agent authentication is only supported for SSH repositories")
		}
	}
	slog.Info("using anonymous authentication")
	return nil, nil
}
