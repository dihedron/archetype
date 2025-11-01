package settings

import "github.com/dihedron/rawdata"

// Auth provides the authentication settings to access a repository.
// It supports authentication via token, username/password, SSH key,
// the default SSH key, or an SSH agent.
type Auth struct {
	Token            *string `short:"T" long:"token" description:"The personal access token for authentication" optional:"true"`
	Username         *string `short:"U" long:"username" description:"The username for authentication" optional:"true"`
	Password         *string `short:"P" long:"password" description:"The password for authentication" optional:"true"`
	SSHKey           *string `short:"K" long:"sshkey" description:"The SSH key for authentication" optional:"true"`
	UseDefaultSSHKey bool    `short:"D" long:"with-default-ssh-key" description:"Use default SSH key for authentication" optional:"true"`
	UseSSHAgent      bool    `short:"A" long:"with-ssh-agent" description:"Use SSH agent for authentication" optional:"true"`
}

// Parameter represents a configurable parameter, including its description,
// type, and default value; these values are used during the post-processing
// of the raw template files to inject the final values into the templates.
type Parameter struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Default     any    `json:"default,omitempty" yaml:"default,omitempty"`
}

// Metadata represents the archetype metadata, which includes the version
// of the metadata structure itself and the set of available parameters.
type Metadata struct {
	Version    int                  `json:"version,omitempty" yaml:"version,omitempty"`
	Parameters map[string]Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// Settings represents the user-provided settings, including the version
// of the settings structure itself and the set of values for the parameters.
type Settings struct {
	Version    int            `json:"version,omitempty" yaml:"version,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// UnmarshalFlag unmarshals a string value into the Settings struct.
// This method is used by the go-flags package to handle custom flag types.
// It takes a string value, which is expected to be in a format that can be
// unmarshalled by the rawdata.UnmarshalInto function (e.g., JSON, YAML), and
// populates the fields of the Settings struct accordingly.
func (s *Settings) UnmarshalFlag(value string) error {
	return rawdata.UnmarshalInto(value, s)
}
