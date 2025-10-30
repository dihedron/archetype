package settings

import "github.com/dihedron/rawdata"

// Auth represents authentication settings for a repository.
type Auth struct {
	Token            *string `short:"T" long:"token" description:"The personal access token for authentication" optional:"true"`
	Username         *string `short:"U" long:"username" description:"The username for authentication" optional:"true"`
	Password         *string `short:"P" long:"password" description:"The password for authentication" optional:"true"`
	SSHKey           *string `short:"K" long:"sshkey" description:"The SSH key for authentication" optional:"true"`
	UseDefaultSSHKey bool    `short:"D" long:"with-default-ssh-key" description:"Use default SSH key for authentication" optional:"true"`
	UseSSHAgent      bool    `short:"A" long:"with-ssh-agent" description:"Use SSH agent for authentication" optional:"true"`
}

// Parameter represents a configurable parameter with its
// name, description, type, default value, and actual value.
// these values will be used when post-processing the raw
// template files.
type Parameter struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Default     any    `json:"default,omitempty" yaml:"default,omitempty"`
}

// Metadata represents the archetype metadata structure,
// including version and parameters.
type Metadata struct {
	Version    int                  `json:"version,omitempty" yaml:"version,omitempty"`
	Parameters map[string]Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// Settings represents the overall settings structure,
// including version and parameters.
type Settings struct {
	Version    int            `json:"version,omitempty" yaml:"version,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// UnmarshalFlag unmarshals a string value into the Settings
// struct; it is used when parsing command-line flags with
// the github.com/jessevdk/go-flags package.
func (s *Settings) UnmarshalFlag(value string) error {
	return rawdata.UnmarshalInto(value, s)
}
