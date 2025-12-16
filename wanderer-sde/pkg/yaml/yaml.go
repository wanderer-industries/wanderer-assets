// Package yaml provides YAML parsing utilities for the SDE converter.
package yaml

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// ParseFile reads and parses a YAML file into the provided target.
func ParseFile(path string, target interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	return Parse(f, target)
}

// Parse decodes YAML from a reader into the provided target.
func Parse(r io.Reader, target interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}
	return nil
}

// ParseFileMap reads and parses a YAML file where the top level is a map.
// This is useful for files like typeIDs.yaml where keys are IDs.
func ParseFileMap[K comparable, V any](path string) (map[K]V, error) {
	var result map[K]V
	if err := ParseFile(path, &result); err != nil {
		return nil, err
	}
	return result, nil
}
