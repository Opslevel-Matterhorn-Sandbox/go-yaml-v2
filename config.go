package main

import (
        "fmt"

        "gopkg.in/yaml.v3"
)

// MapSlice encodes and decodes as a YAML map.
// The order of keys is preserved when encoding and decoding.
type MapSlice []MapItem

// MapItem is an item in a MapSlice.
type MapItem struct {
        Key, Value interface{}
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (ms *MapSlice) UnmarshalYAML(value *yaml.Node) error {
        if value.Kind != yaml.MappingNode {
                return fmt.Errorf("expected a mapping node")
        }
        *ms = make(MapSlice, 0, len(value.Content)/2)
        for i := 0; i < len(value.Content); i += 2 {
                var key, val interface{}
                if err := value.Content[i].Decode(&key); err != nil {
                        return err
                }
                if err := value.Content[i+1].Decode(&val); err != nil {
                        return err
                }
                *ms = append(*ms, MapItem{Key: key, Value: val})
        }
        return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (ms MapSlice) MarshalYAML() (interface{}, error) {
        node := &yaml.Node{
                Kind: yaml.MappingNode,
        }
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
