package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func LoadConfig(data []byte) (yaml.MapSlice, error) {
	var out yaml.MapSlice
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return out, nil
}

func SaveConfig(c yaml.MapSlice) ([]byte, error) {
	return yaml.Marshal(c)
}

func Get(c yaml.MapSlice, key string) (interface{}, bool) {
	for _, item := range c {
		if k, ok := item.Key.(string); ok && k == key {
			return item.Value, true
		}
	}
	return nil, false
}

func Merge(base, overrides yaml.MapSlice) yaml.MapSlice {
	seen := make(map[string]bool)
	var out yaml.MapSlice
	for _, item := range base {
		k, ok := item.Key.(string)
		if !ok {
			out = append(out, item)
			continue
		}
		seen[k] = true
		if v, found := getKey(overrides, k); found {
			out = append(out, yaml.MapItem{Key: k, Value: v})
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

func getKey(c yaml.MapSlice, key string) (interface{}, bool) {
	for _, item := range c {
		if k, ok := item.Key.(string); ok && k == key {
			return item.Value, true
		}
	}
	return nil, false
}
