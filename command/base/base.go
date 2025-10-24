package base

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/dihedron/archetype/repository"
)

// This is the set of common command line options.
type Command struct {
	Repository string `short:"r" long:"repository" description:"The Git repository to clone" required:"true"`
	Tag        string `short:"t" long:"tag" description:"The tag or commit to clone" optional:"true" default:"latest"`
	Auth       struct {
		Token             *string `short:"T" long:"token" description:"The personal access token for authentication" optional:"true"`
		Username          *string `short:"U" long:"username" description:"The username for authentication" optional:"true"`
		Password          *string `short:"P" long:"password" description:"The password for authentication" optional:"true"`
		SSHKey            *string `short:"K" long:"sshkey" description:"The SSH key for authentication" optional:"true"`
		WithDefaultSSHKey bool    `short:"D" long:"with-default-ssh-key" description:"Use default SSH key for authentication" optional:"true"`
		WithSSHAgent      bool    `short:"A" long:"with-ssh-agent" description:"Use SSH agent for authentication" optional:"true"`
	} `group:"Authentication Options" description:"Authentication options"`
}

// AuthenticationOpts returns the repository.Option needed
// to configure authenticated requests against the remote repository.
func (cmd *Command) AuthenticationOpts() (repository.Option, error) {
	if cmd.Auth.Token != nil {
		if strings.HasPrefix(cmd.Repository, "http") {
			slog.Info("using token for authentication")
			return repository.WithTokenAuth(*cmd.Auth.Token), nil
		} else {
			slog.Error("token authentication is only supported for HTTP repositories")
			return nil, errors.New("token authentication is only supported for HTTP repositories")
		}
	} else if cmd.Auth.Password != nil && cmd.Auth.Username != nil {
		if strings.HasPrefix(cmd.Repository, "http") {
			slog.Info("using username and password for authentication")
			return repository.WithBasicAuth(*cmd.Auth.Username, *cmd.Auth.Password), nil
		} else {
			slog.Error("username and password authentication is only supported for HTTP repositories")
			return nil, errors.New("username and password authentication is only supported for HTTP repositories")
		}
	} else if cmd.Auth.SSHKey != nil {
		if strings.HasPrefix(cmd.Repository, "ssh://") || strings.HasPrefix(cmd.Repository, "git@") {
			slog.Info("using SSH key for authentication")
			return repository.WithSSHKey(*cmd.Auth.SSHKey, nil), nil
		} else {
			slog.Error("SSH key authentication is only supported for SSH repositories")
			return nil, errors.New("SSH key authentication is only supported for SSH repositories")
		}
	} else if cmd.Auth.WithDefaultSSHKey {
		if strings.HasPrefix(cmd.Repository, "ssh://") || strings.HasPrefix(cmd.Repository, "git@") {
			slog.Info("using default SSH key for authentication")
			return repository.WithDefaultSSHKey(), nil
		} else {
			slog.Error("SSH key authentication is only supported for SSH repositories")
			return nil, errors.New("SSH key authentication is only supported for SSH repositories")
		}
	} else if cmd.Auth.WithSSHAgent {
		if strings.HasPrefix(cmd.Repository, "ssh://") || strings.HasPrefix(cmd.Repository, "git@") {
			slog.Info("using SSH agent for authentication")
			return repository.WithSSHAgent(), nil
		} else {
			slog.Error("SSH agent authentication is only supported for SSH repositories")
			return nil, errors.New("SSH agent authentication is only supported for SSH repositories")
		}
	}
	slog.Info("using anonymous authentication")
	return nil, nil
}
