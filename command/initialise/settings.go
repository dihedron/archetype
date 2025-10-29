package initialise

/*
type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Default     string `json:"default,omitempty" yaml:"default,omitempty"`
	Value       any    `json:"value,omitempty" yaml:"value,omitempty"`
}

type Auth struct {
	Token            *string `json:"token,omitempty" yaml:"token,omitempty" short:"T" long:"token" description:"The personal access token for authentication" optional:"true"`
	Username         *string `json:"username,omitempty" yaml:"username,omitempty" short:"U" long:"username" description:"The username for authentication" optional:"true"`
	Password         *string `json:"password,omitempty" yaml:"password,omitempty" short:"P" long:"password" description:"The password for authentication" optional:"true"`
	SSHKey           *string `json:"sshKey,omitempty" yaml:"sshKey,omitempty" short:"K" long:"sshkey" description:"The SSH key for authentication" optional:"true"`
	UseDefaultSSHKey bool    `json:"useDefaultSSHKey,omitempty" yaml:"useDefaultSSHKey,omitempty" short:"D" long:"with-default-ssh-key" description:"Use default SSH key for authentication" optional:"true"`
	UseSSHAgent      bool    `json:"useSSHAgent,omitempty" yaml:"useSSHAgent,omitempty" short:"A" long:"with-ssh-agent" description:"Use SSH agent for authentication" optional:"true"`
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
*/
