package settings

type Setting struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Value       any    `json:"value" yaml:"value"`
}

type Settings struct {
	Values []Setting
}
