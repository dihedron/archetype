package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Type int8

const (
	String Type = iota
	Number
	Boolean
)

// String converts a Type into its string representation.
func (t Type) String() string {
	return []string{
		"string",
		"number",
		"boolean",
	}[t]
}

// Parse parses a string into a Type value.
func (t *Type) Parse(value string) error {
	switch strings.ToLower(value) {
	case "string":
		*t = String
	case "number":
		*t = Number
	case "boolean":
		*t = Boolean
	default:
		return errors.New("unsupported Type value")
	}
	return nil
}

// MarshalYAML marshals the Type value into a YAML string.
func (t Type) MarshalYAML() (interface{}, error) {
	v := t.String()
	if v != "" {
		return v, nil
	}
	return "", fmt.Errorf("unsupported Type value: %s", v)
}

// UnmarshalYAML unmarshals a YAML value into a Type value.
func (t *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	size := ""
	err := unmarshal(&size)
	if err != nil {
		return err
	}
	return t.Parse(size)
}

// MarshalJSON marshals the Type value into a JSON string.
func (t Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON unmarshals a JSON value into a Type value.
func (t *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	t.Parse(value)
	return nil
}
