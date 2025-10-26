package initialise

import (
	"github.com/dihedron/rawdata"
)

type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Default     string `json:"default,omitempty" yaml:"default,omitempty"`
	Value       any    `json:"value,omitempty" yaml:"value,omitempty"`
}

type Auth struct {
	Token    *string `json:"token,omitempty" yaml:"token,omitempty"`
	Username *string `json:"username,omitempty" yaml:"username,omitempty"`
	Password *string `json:"password,omitempty" yaml:"password,omitempty"`
	SSHKey   *string `json:"sshKey,omitempty" yaml:"sshKey,omitempty"`
}

type Repository struct {
	URL  string  `json:"url,omitempty" yaml:"url,omitempty"`
	Tag  *string `json:"tag,omitempty" yaml:"tag,omitempty"`
	Auth *Auth   `json:"auth,omitempty" yaml:"auth,omitempty"`
}

type Settings struct {
	Version    int         `json:"version,omitempty" yaml:"version,omitempty"`
	Repository Repository  `json:"repository" yaml:"repository"`
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func (s *Settings) UnmarshalFlag(value string) error {
	return rawdata.UnmarshalInto(value, s)
}
