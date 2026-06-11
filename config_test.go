package main

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// makeMapping is a helper that builds a Config (MappingNode) from alternating
// string key / interface{} value pairs.
func makeMapping(pairs ...interface{}) Config {
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for i := 0; i+1 < len(pairs); i += 2 {
		key := pairs[i].(string)
		val := pairs[i+1]

		keyNode := &yaml.Node{}
		if err := keyNode.Encode(key); err != nil {
			panic(err)
		}
		// Encode wraps in a DocumentNode; unwrap it.
		if keyNode.Kind == yaml.DocumentNode {
			keyNode = keyNode.Content[0]
		}

		valNode := &yaml.Node{}
		if err := valNode.Encode(val); err != nil {
			panic(err)
		}
		if valNode.Kind == yaml.DocumentNode {
			valNode = valNode.Content[0]
		}

		node.Content = append(node.Content, keyNode, valNode)
	}
	return node
}

func TestLoadConfig(t *testing.T) {
	data := []byte("a: 1\nb: 2\nc: 3")
	cfg, err := LoadConfig(data)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	// A mapping node with 3 key-value pairs has 6 content entries.
	if len(cfg.Content) != 6 {
		t.Errorf("expected 6 content nodes (3 pairs), got %d", len(cfg.Content))
	}
	keys := []string{"a", "b", "c"}
	for i, want := range keys {
		if cfg.Content[i*2].Value != want {
			t.Errorf("key[%d]: got %q, want %q", i, cfg.Content[i*2].Value, want)
		}
	}
}

func TestSaveConfig(t *testing.T) {
	cfg := makeMapping("name", "test", "version", "2.0")
	out, err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	roundTrip, err := LoadConfig(out)
	if err != nil {
		t.Fatalf("LoadConfig round-trip: %v", err)
	}
	// Expect 2 key-value pairs = 4 content nodes.
	if len(roundTrip.Content) != 4 {
		t.Errorf("round-trip length: got %d content nodes, want 4", len(roundTrip.Content))
	}
}

func TestGet(t *testing.T) {
	cfg := makeMapping("name", "myapp", "count", 42)
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
	base := makeMapping("name", "app", "env", "dev", "replicas", 1)
	ov := makeMapping("env", "prod", "replicas", 5)
	merged := Merge(base, ov)
	// 3 key-value pairs = 6 content nodes.
	if len(merged.Content) != 6 {
		t.Fatalf("merged content length: got %d, want 6", len(merged.Content))
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
	// Order: name, env, replicas
	if merged.Content[0].Value != "name" || merged.Content[2].Value != "env" || merged.Content[4].Value != "replicas" {
		t.Errorf("merge order: keys are %q, %q, %q",
			merged.Content[0].Value, merged.Content[2].Value, merged.Content[4].Value)
	}
}

func TestMergeAddNewKeys(t *testing.T) {
	base := makeMapping("a", 1)
	ov := makeMapping("b", 2, "c", 3)
	merged := Merge(base, ov)
	// 3 key-value pairs = 6 content nodes.
	if len(merged.Content) != 6 {
		t.Fatalf("merged content length: got %d, want 6", len(merged.Content))
	}
	keys := []string{"a", "b", "c"}
	for i, want := range keys {
		if merged.Content[i*2].Value != want {
			t.Errorf("key[%d]: got %q, want %q", i, merged.Content[i*2].Value, want)
		}
	}
}

func TestMapItemPreservesOrder(t *testing.T) {
	cfg := makeMapping("first", 1, "second", 2)
	out, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	back, err := LoadConfig(out)
	if err != nil {
		t.Fatal(err)
	}
	if back.Content[0].Value != "first" || back.Content[2].Value != "second" {
		t.Errorf("order changed: %q, %q", back.Content[0].Value, back.Content[2].Value)
	}
}
