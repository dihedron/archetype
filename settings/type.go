package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Type represents the type of a parameter.
type Type int8

const (
	String Type = iota
	Number
	Boolean
)

// String returns the string representation of the Type.
func (t Type) String() string {
	return []string{
		"string",
		"number",
		"boolean",
	}[t]
}

// Parse parses a string and sets the Type value accordingly.
func (t *Type) Parse(value string) error {
	switch strings.ToLower(value) {
	case "string":
		*t = String
	case "number", "int", "integer":
		*t = Number
	case "boolean", "bool":
		*t = Boolean
	default:
		return errors.New("unsupported Type value")
	}
	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (t Type) MarshalYAML() (interface{}, error) {
	v := t.String()
	if v != "" {
		return v, nil
	}
	return "", fmt.Errorf("unsupported Type value: %s", v)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (t *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	size := ""
	err := unmarshal(&size)
	if err != nil {
		return err
	}
	return t.Parse(size)
}

// MarshalJSON implements the json.Marshaler interface.
func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	t.Parse(value)
	return nil
}
