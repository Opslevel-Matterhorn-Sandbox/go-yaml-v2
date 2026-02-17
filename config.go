package main

import (
        "fmt"

        "gopkg.in/yaml.v3"
)

// MapItem holds a single key-value pair from an ordered YAML mapping.
type MapItem struct {
        Key, Value interface{}
}

// MapSlice is an ordered sequence of key-value pairs that round-trips through
// YAML while preserving the original key order.
type MapSlice []MapItem

// MarshalYAML encodes the MapSlice as an ordered YAML mapping node so that
// key order is preserved when marshaling with gopkg.in/yaml.v3.
func (ms MapSlice) MarshalYAML() (interface{}, error) {
        node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
        for _, item := range ms {
                keyNode := &yaml.Node{}
                if err := keyNode.Encode(item.Key); err != nil {
                        return nil, err
                }
                valNode := &yaml.Node{}
                if err := valNode.Encode(item.Value); err != nil {
                        return nil, err
                }
                node.Content = append(node.Content, keyNode, valNode)
        }
        return node, nil
}

// UnmarshalYAML decodes an ordered YAML mapping node into the MapSlice so
// that key order is preserved when unmarshaling with gopkg.in/yaml.v3.
func (ms *MapSlice) UnmarshalYAML(value *yaml.Node) error {
        if value.Kind == yaml.DocumentNode {
                if len(value.Content) == 0 {
                        return nil
                }
                return ms.UnmarshalYAML(value.Content[0])
        }
        if value.Kind != yaml.MappingNode {
                return fmt.Errorf("unmarshal config: expected a mapping node, got %v", value.Kind)
        }
        for i := 0; i+1 < len(value.Content); i += 2 {
                var k interface{}
                if err := value.Content[i].Decode(&k); err != nil {
                        return err
                }
                var v interface{}
                if err := value.Content[i+1].Decode(&v); err != nil {
                        return err
                }
                *ms = append(*ms, MapItem{Key: k, Value: v})
        }
        return nil
}

func LoadConfig(data []byte) (MapSlice, error) {
        var out MapSlice
        if err := yaml.Unmarshal(data, &out); err != nil {
                return nil, fmt.Errorf("unmarshal config: %w", err)
        }
        return out, nil
}

func SaveConfig(c MapSlice) ([]byte, error) {
        return yaml.Marshal(c)
}

func Get(c MapSlice, key string) (interface{}, bool) {
        for _, item := range c {
                if k, ok := item.Key.(string); ok && k == key {
                        return item.Value, true
                }
        }
        return nil, false
}

func Merge(base, overrides MapSlice) MapSlice {
        seen := make(map[string]bool)
        var out MapSlice
        for _, item := range base {
                k, ok := item.Key.(string)
                if !ok {
                        out = append(out, item)
                        continue
                }
                seen[k] = true
                if v, found := getKey(overrides, k); found {
                        out = append(out, MapItem{Key: k, Value: v})
                } else {
                        out = append(out, item)
                }
        }
        for _, item := range overrides {
                k, ok := item.Key.(string)
                if !ok || seen[k] {
                        continue
                }
                seen[k] = true
                out = append(out, item)
        }
        return out
}

func getKey(c MapSlice, key string) (interface{}, bool) {
        for _, item := range c {
                if k, ok := item.Key.(string); ok && k == key {
                        return item.Value, true
                }
        }
        return nil, false
}
