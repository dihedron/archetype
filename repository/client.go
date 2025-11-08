package repository

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"github.com/go-git/go-git/v6/storage/memory"
)

// Repository represents a Git repository, either local or remote.
type Repository struct {
	address    string
	auth       transport.AuthMethod
	repository *git.Repository
	proxy      *transport.ProxyOptions
}

// Option is a functional option for configuring a Repository.
type Option func(*Repository)

// New creates a new Repository with the given address and options.
func New(address string, options ...Option) (*Repository, error) {
	if address == "" || address == "." {
		slog.Debug("using default address", "address", "file://./")
		address = "file://./"
	} else {
		slog.Debug("using explicitly provided address", "address", address)
	}
	repository := &Repository{
		address: address,
	}
	for _, option := range options {
		if option != nil {
			option(repository)
		}
	}
	if strings.HasPrefix(repository.address, "file://") {
		err := repository.open()
		if err != nil {
			slog.Error("failed to open local repository", "error", err)
			return nil, err
		}
		return repository, nil
	} else if strings.HasPrefix(repository.address, "ssh://") || strings.HasPrefix(repository.address, "http://") || strings.HasPrefix(repository.address, "https://") {
		err := repository.clone()
		if err != nil {
			slog.Error("failed to clone remote repository", "error", err)
			return nil, err
		}
	}
	return repository, nil
}

// WithBasicAuth configures the Repository to use HTTP basic authentication.
func WithBasicAuth(username string, password string) Option {
	return func(repository *Repository) {
		repository.auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}
}

// WithTokenAuth configures the Repository to use a personal access token for
// authentication.
func WithTokenAuth(token string) Option {
	return func(repository *Repository) {
		repository.auth = &http.BasicAuth{
			Username: "git",
			Password: token,
		}
	}
}

// WithSSHKey configures the Repository to use an SSH key for authentication.
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

// WithDefaultSSHKey configures the Repository to use the default SSH key for
// authentication.
func WithDefaultSSHKey() Option {
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("failed to get user home directory", "error", err)
		os.Exit(1)
	}
	return WithSSHKey(filepath.Join(home, ".ssh", "id_rsa"), nil)
}

// WithSSHAgent configures the Repository to use an SSH agent for authentication.
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

// WithProxy configures the Repository to use a proxy.
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

// WithProxyFromEnv configures the Repository to use a proxy from the environment
// variables (HTTP_PROXY, HTTPS_PROXY).
func WithProxyFromEnv() Option {
	return func(repository *Repository) {
		slog.Info("setting up proxy for git transport", "repository", repository.address)
		var (
			proxyURL string
			ok       bool
		)
		if strings.HasPrefix(repository.address, "https") {
			slog.Debug("retrieving proxy for HTTPs address")
			if proxyURL, ok = os.LookupEnv("HTTPS_PROXY"); !ok {
				proxyURL, ok = os.LookupEnv("https_proxy")
				if !ok {
					slog.Debug("no proxy")
				}
			}
		} else if strings.HasPrefix(repository.address, "http") {
			slog.Debug("retrieving proxy for HTTP address")
			if proxyURL, ok = os.LookupEnv("HTTP_PROXY"); !ok {
				proxyURL, ok = os.LookupEnv("http_proxy")
				if !ok {
					slog.Debug("no proxy")
				}
			}
		}
		if proxyURL == "" {
			slog.Debug("no proxy available in environment")
			return
		}
		slog.Debug("retrieved HTTP(s) proxy URL from environment", "url", proxyURL)
		if parsed, err := url.Parse(proxyURL); err != nil {
			slog.Error("invalid proxy URL", "error", err)
		} else {
			var (
				username string
				password string
				address  string
			)
			if parsed.User != nil {
				username = parsed.User.Username()
				if password, ok = parsed.User.Password(); !ok {
					username = ""
				}
			}
			parsed.User = nil
			address = parsed.String()
			slog.Debug("setting auto proxy", "url", address, "username", username, "password", password)
			repository.proxy = &transport.ProxyOptions{
				URL:      address,
				Username: username,
				Password: password,
			}
		}
	}
}

// open opens the repository using the povided address.
func (r *Repository) open() error {
	if r.address == "" {
		slog.Error("repository address not set")
		return fmt.Errorf("repository address not set")
	}
	directory := strings.TrimPrefix(r.address, "file://")
	slog.Debug("opening repository", "address", r.address, "directory", directory)
	repository, err := git.PlainOpen(directory)
	if err != nil {
		slog.Error("failed to open repository", "error", err)
		return err
	}
	r.repository = repository
	return nil
}

// Clone clones the repository into memory.
func (r *Repository) clone() error {
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
