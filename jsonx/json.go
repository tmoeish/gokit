// Package jsonx provides JSON and YAML utility functions.
package jsonx

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Marshal returns the JSON encoding of v, or "" on error.
func Marshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// MarshalPretty returns indented JSON encoding of v, or "" on error.
func MarshalPretty(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

// MarshalBytes returns the JSON encoding of v as bytes.
// Returns nil on error.
func MarshalBytes(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return b
}

// MustMarshal returns the JSON encoding of v, panicking on error.
func MustMarshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Unmarshal decodes JSON data into v.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// UnmarshalStr decodes a JSON string into v.
func UnmarshalStr(s string, v any) error {
	return json.Unmarshal([]byte(s), v)
}

// MustUnmarshal decodes JSON data into v, panicking on error.
func MustUnmarshal(data []byte, v any) {
	if err := json.Unmarshal(data, v); err != nil {
		panic(err)
	}
}

// Clone deep-clones v via JSON round-trip into a new value of the same type.
func Clone[T any](v T) (T, error) {
	b, err := json.Marshal(v)
	if err != nil {
		var zero T
		return zero, err
	}
	var result T
	if err := json.Unmarshal(b, &result); err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

// MarshalYAML returns the YAML encoding of v, or "" on error.
func MarshalYAML(v any) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// MarshalYAMLBytes returns the YAML encoding of v as bytes.
func MarshalYAMLBytes(v any) []byte {
	b, err := yaml.Marshal(v)
	if err != nil {
		return nil
	}
	return b
}

// UnmarshalYAML decodes YAML data into v.
func UnmarshalYAML(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

// UnmarshalYAMLStr decodes a YAML string into v.
func UnmarshalYAMLStr(s string, v any) error {
	return yaml.Unmarshal([]byte(s), v)
}

// JSONToYAML converts a JSON string to YAML.
func JSONToYAML(jsonStr string) (string, error) {
	var obj any
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return "", err
	}
	return MarshalYAML(obj), nil
}

// YAMLToJSON converts a YAML string to JSON.
func YAMLToJSON(yamlStr string) (string, error) {
	var obj any
	if err := yaml.Unmarshal([]byte(yamlStr), &obj); err != nil {
		return "", err
	}
	return Marshal(obj), nil
}

// Get retrieves a nested value from a JSON object by dot-separated path.
// Returns nil if the path does not exist.
func Get(data []byte, path ...string) any {
	var obj any
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil
	}
	return getPath(obj, path)
}

func getPath(obj any, path []string) any {
	if len(path) == 0 {
		return obj
	}
	m, ok := obj.(map[string]any)
	if !ok {
		return nil
	}
	return getPath(m[path[0]], path[1:])
}
