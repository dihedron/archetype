package parameters

import "github.com/dihedron/rawdata"

type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Default     string `json:"default,omitempty" yaml:"default,omitempty"`
	Value       any    `json:"value,omitempty" yaml:"value,omitempty"`
}

type Parameters struct {
	Version int         `json:"version,omitempty" yaml:"version,omitempty"`
	Values  []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func (p *Parameters) UnmarshalFlag(value string) error {
	return rawdata.UnmarshalInto(value, p)
}
