package main

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadConfig(t *testing.T) {
	data := []byte("a: 1\nb: 2\nc: 3")
	cfg, err := LoadConfig(data)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Kind != yaml.MappingNode {
		t.Fatalf("expected MappingNode, got %v", cfg.Kind)
	}
	// Content has 6 nodes: 3 keys + 3 values
	if len(cfg.Content) != 6 {
		t.Errorf("expected 6 nodes (3 key-value pairs), got %d", len(cfg.Content))
	}
	if cfg.Content[0].Value != "a" || cfg.Content[2].Value != "b" || cfg.Content[4].Value != "c" {
		t.Errorf("key order not preserved")
	}
}

func TestSaveConfig(t *testing.T) {
	cfg := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "name"},
			{Kind: yaml.ScalarNode, Value: "test"},
			{Kind: yaml.ScalarNode, Value: "version"},
			{Kind: yaml.ScalarNode, Value: "2.0"},
		},
	}
	out, err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	roundTrip, err := LoadConfig(out)
	if err != nil {
		t.Fatalf("LoadConfig round-trip: %v", err)
	}
	if len(roundTrip.Content) != 4 {
		t.Errorf("round-trip length: got %d nodes (expected 4)", len(roundTrip.Content))
	}
}

func TestGet(t *testing.T) {
	cfg := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "name"},
			{Kind: yaml.ScalarNode, Value: "myapp"},
			{Kind: yaml.ScalarNode, Value: "count"},
			{Kind: yaml.ScalarNode, Value: "42"},
		},
	}
	v, ok := Get(cfg, "name")
	if !ok || v != "myapp" {
		t.Errorf("Get(name): got %v, %v", v, ok)
	}
	v, ok = Get(cfg, "count")
	if !ok || v != 42 {
		t.Errorf("Get(count): got %v, %v", v, ok)
	}
	_, ok = Get(cfg, "missing")
	if ok {
		t.Error("Get(missing): expected false")
	}
}

func TestMerge(t *testing.T) {
	base := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "name"},
			{Kind: yaml.ScalarNode, Value: "app"},
			{Kind: yaml.ScalarNode, Value: "env"},
			{Kind: yaml.ScalarNode, Value: "dev"},
			{Kind: yaml.ScalarNode, Value: "replicas"},
			{Kind: yaml.ScalarNode, Value: "1"},
		},
	}
	ov := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "env"},
			{Kind: yaml.ScalarNode, Value: "prod"},
			{Kind: yaml.ScalarNode, Value: "replicas"},
			{Kind: yaml.ScalarNode, Value: "5"},
		},
	}
	merged := Merge(base, ov)
	if len(merged.Content) != 6 {
		t.Fatalf("merged length: got %d nodes (expected 6)", len(merged.Content))
	}
	if v, _ := Get(merged, "name"); v != "app" {
		t.Errorf("name: got %v", v)
	}
	if v, _ := Get(merged, "env"); v != "prod" {
		t.Errorf("env: got %v", v)
	}
	if v, _ := Get(merged, "replicas"); v != 5 {
		t.Errorf("replicas: got %v", v)
	}
	if merged.Content[0].Value != "name" || merged.Content[2].Value != "env" || merged.Content[4].Value != "replicas" {
		t.Errorf("merge order not preserved")
	}
}

func TestMergeAddNewKeys(t *testing.T) {
	base := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "a"},
			{Kind: yaml.ScalarNode, Value: "1"},
		},
	}
	ov := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "b"},
			{Kind: yaml.ScalarNode, Value: "2"},
			{Kind: yaml.ScalarNode, Value: "c"},
			{Kind: yaml.ScalarNode, Value: "3"},
		},
	}
	merged := Merge(base, ov)
	if len(merged.Content) != 6 {
		t.Fatalf("merged length: got %d nodes (expected 6)", len(merged.Content))
	}
	if merged.Content[0].Value != "a" || merged.Content[2].Value != "b" || merged.Content[4].Value != "c" {
		t.Errorf("merge order not preserved")
	}
}

func TestMapItemPreservesOrder(t *testing.T) {
	node := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "first"},
			{Kind: yaml.ScalarNode, Value: "1"},
			{Kind: yaml.ScalarNode, Value: "second"},
			{Kind: yaml.ScalarNode, Value: "2"},
		},
	}
	out, err := yaml.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	back, err := LoadConfig(out)
	if err != nil {
		t.Fatal(err)
	}
	if back.Content[0].Value != "first" || back.Content[2].Value != "second" {
		t.Errorf("order changed")
	}
}
