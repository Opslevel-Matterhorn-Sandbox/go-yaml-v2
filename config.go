package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Config is an order-preserving YAML mapping represented as a yaml.Node.
type Config = *yaml.Node

// LoadConfig parses YAML data and returns the top-level mapping node.
func LoadConfig(data []byte) (Config, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		return doc.Content[0], nil
	}
	// Return an empty mapping node when the input is empty.
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}, nil
}

// SaveConfig marshals the mapping node back to YAML bytes.
func SaveConfig(c Config) ([]byte, error) {
	return yaml.Marshal(c)
}

// Get returns the value associated with key in the mapping node, preserving
// insertion order. The node Content holds alternating key/value pairs.
func Get(c Config, key string) (interface{}, bool) {
	if c == nil || c.Kind != yaml.MappingNode {
		return nil, false
	}
	for i := 0; i+1 < len(c.Content); i += 2 {
		if c.Content[i].Value == key {
			var v interface{}
			if err := c.Content[i+1].Decode(&v); err == nil {
				return v, true
			}
		}
	}
	return nil, false
}

// Merge produces a new Config that starts with all keys from base, with any
// values overridden by overrides, and then appends any extra keys from
// overrides that were not present in base. Key order is preserved.
func Merge(base, overrides Config) Config {
	out := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}

	seen := make(map[string]bool)

	// Walk base entries; override value when the key exists in overrides.
	for i := 0; i+1 < len(base.Content); i += 2 {
		keyNode := base.Content[i]
		k := keyNode.Value
		seen[k] = true
		if ov := findValue(overrides, k); ov != nil {
			out.Content = append(out.Content, keyNode, ov)
		} else {
			out.Content = append(out.Content, keyNode, base.Content[i+1])
		}
	}

	// Append keys from overrides that were not in base.
	for i := 0; i+1 < len(overrides.Content); i += 2 {
		k := overrides.Content[i].Value
		if !seen[k] {
			seen[k] = true
			out.Content = append(out.Content, overrides.Content[i], overrides.Content[i+1])
		}
	}

	return out
}

// findValue returns the value node for key in a mapping node, or nil.
func findValue(c Config, key string) *yaml.Node {
	if c == nil || c.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(c.Content); i += 2 {
		if c.Content[i].Value == key {
			return c.Content[i+1]
		}
	}
	return nil
}
