package repository

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"github.com/go-git/go-git/v6/storage/memory"
)

// Repository is a struct that encapsulates the configuration and methods needed to
// clone a Git repository. It supports various authentication methods and allows
// specifying branch, commit, and target directory for the clone operation.
type Repository struct {
	address    string
	auth       transport.AuthMethod
	repository *git.Repository
	proxy      *transport.ProxyOptions
}

// Option is a functional option for configuring the Cloner.
type Option func(*Repository)

// New creates a new Cloner instance with the specified repository URL and
// applies any provided functional options to configure it.
func New(address string, options ...Option) *Repository {
	repository := &Repository{
		address: address,
	}
	for _, option := range options {
		if option != nil {
			option(repository)
		}
	}
	return repository
}

// WithBasicAuth is a helper to set up HTTP basic authentication using a
// username and password.
func WithBasicAuth(username string, password string) Option {
	return func(repository *Repository) {
		repository.auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}
}

// WithTokenAuth is a helper to set up HTTP Basic authentication using a
// personal access token (PAT).
func WithTokenAuth(token string) Option {
	return func(repository *Repository) {
		repository.auth = &http.BasicAuth{
			Username: "git",
			Password: token,
		}
	}
}

// WithSSHKey is a helper to set up SSH authentication using a private key
// file. If the key is password-protected, the password must be provided as
// well; otherwise, nil can be passed.
func WithSSHKey(path string, password *string) Option {
	return func(repository *Repository) {
		slog.Info("setting up SSH authentication...")
		pwd := ""
		if password != nil {
			pwd = *password
		}
		keys, err := ssh.NewPublicKeysFromFile("git", path, pwd)
		if err != nil {
			slog.Error("failed to generate public keys from file", "error", err)
			os.Exit(1)
		}
		// determine the known hosts dynamically from the user's home directory
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("failed to get user home directory", "error", err)
			os.Exit(1)
		}
		callback, err := ssh.NewKnownHostsCallback(filepath.Join(home, ".ssh", "known_hosts"))
		if err != nil {
			slog.Error("failed to create known_hosts callback", "error", err)
			os.Exit(1)
		}
		keys.HostKeyCallback = callback
		repository.auth = keys
	}
}

// WithDefaultSSHKey is a helper to set up SSH authentication using the default
// private key file located in the user's home directory (~/.ssh/id_rsa).
func WithDefaultSSHKey() Option {
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("failed to get user home directory", "error", err)
		os.Exit(1)
	}
	return WithSSHKey(filepath.Join(home, ".ssh", "id_rsa"), nil)
}

// WithSSHAgent is a helper to set up SSH authentication using the SSH agent
// (ssh-agent) running on the local machine.
func WithSSHAgent() Option {
	return func(repository *Repository) {
		slog.Info("setting up SSH authentication using SSH agent...")
		authMethod, err := ssh.NewSSHAgentAuth("git")
		if err != nil {
			slog.Error("failed to connect to SSH agent", "error", err)
			os.Exit(1)
		}
		repository.auth = authMethod
	}
}

// WithProxy is a helper to set up a proxy for HTTP/HTTPS transport.
func WithProxy(proxyURL string, username string, password string) Option {
	return func(repository *Repository) {
		slog.Info("setting up proxy for git transport", "proxy", proxyURL)
		repository.proxy = &transport.ProxyOptions{
			URL:      proxyURL,
			Username: username,
			Password: password,
		}
	}
}

// Clone performs the clone operation using the configured settings and
// authentication method. It clones the repository into memory.
func (r *Repository) Clone() error {
	slog.Info("creating in-memory storage...")
	storage := memory.NewStorage()

	slog.Info("cloning repository into memory", "address", r.address)

	options := &git.CloneOptions{
		URL:      r.address,
		Auth:     r.auth,
		Progress: os.Stdout,
	}
	if r.proxy != nil && strings.HasPrefix(r.proxy.URL, "http") {
		slog.Debug("applying proxy settings", "address", r.address, "proxy", r.proxy.URL)
		options.ProxyOptions = *r.proxy
	}

	repository, err := git.Clone(storage, nil, options)
	if err != nil {
		slog.Error("failed to clone repository", "error", err)
		return err
	}
	slog.Info("clone successful!")

	head, err := repository.Head()
	if err != nil {
		slog.Error("failed to get repository HEAD", "error", err)
		return err
	}
	slog.Debug("repository HEAD found", "commit", head.Hash())
	r.repository = repository
	return nil
}
