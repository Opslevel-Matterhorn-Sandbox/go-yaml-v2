package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func LoadConfig(data []byte) (*yaml.Node, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	// The unmarshalled node is a document node; we want the first content (the map)
	if len(doc.Content) == 0 {
		return nil, fmt.Errorf("empty document")
	}
	return doc.Content[0], nil
}

func SaveConfig(c *yaml.Node) ([]byte, error) {
	return yaml.Marshal(c)
}

func Get(c *yaml.Node, key string) (interface{}, bool) {
	if c.Kind != yaml.MappingNode {
		return nil, false
	}
	for i := 0; i < len(c.Content); i += 2 {
		if i+1 < len(c.Content) && c.Content[i].Value == key {
			// Decode the value node into an interface{}
			var val interface{}
			if err := c.Content[i+1].Decode(&val); err == nil {
				return val, true
			}
			return nil, false
		}
	}
	return nil, false
}

func Merge(base, overrides *yaml.Node) *yaml.Node {
	if base.Kind != yaml.MappingNode || overrides.Kind != yaml.MappingNode {
		return base
	}

	seen := make(map[string]bool)
	var content []*yaml.Node

	// Process base nodes
	for i := 0; i < len(base.Content); i += 2 {
		if i+1 >= len(base.Content) {
			break
		}
		keyNode := base.Content[i]
		valueNode := base.Content[i+1]
		k := keyNode.Value

		seen[k] = true
		if overrideValue := getKeyNode(overrides, k); overrideValue != nil {
			content = append(content, keyNode, overrideValue)
		} else {
			content = append(content, keyNode, valueNode)
		}
	}

	// Add new keys from overrides
	for i := 0; i < len(overrides.Content); i += 2 {
		if i+1 >= len(overrides.Content) {
			break
		}
		keyNode := overrides.Content[i]
		valueNode := overrides.Content[i+1]
		k := keyNode.Value

		if !seen[k] {
			seen[k] = true
			content = append(content, keyNode, valueNode)
		}
	}

	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: content,
	}
}

func getKeyNode(c *yaml.Node, key string) *yaml.Node {
	if c.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(c.Content); i += 2 {
		if i+1 < len(c.Content) && c.Content[i].Value == key {
			return c.Content[i+1]
		}
	}
	return nil
}
