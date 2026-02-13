package main

import (
        "testing"
)

func TestEmptyYAML(t *testing.T) {
        empty := []byte("")
        cfg, err := LoadConfig(empty)
        if err != nil {
                t.Fatalf("Empty YAML should parse: %v", err)
        }
        if len(cfg) != 0 {
                t.Errorf("Empty YAML should result in empty slice")
        }
}

func TestSingleKeyValue(t *testing.T) {
        single := []byte("key: value")
        cfg, err := LoadConfig(single)
        if err != nil {
                t.Fatalf("Single key-value failed: %v", err)
        }
        if len(cfg) != 1 {
                t.Errorf("Expected 1 item, got %d", len(cfg))
        }
}

func TestComplexValues(t *testing.T) {
        complex := []byte("int: 42\nfloat: 3.14\nbool: true\nstring: hello")
        cfg, err := LoadConfig(complex)
        if err != nil {
                t.Fatalf("Complex values failed: %v", err)
        }
        if len(cfg) != 4 {
                t.Errorf("Expected 4 items, got %d", len(cfg))
        }
}

func TestRoundTripPreservation(t *testing.T) {
        original := []byte("first: 1\nsecond: 2\nthird: 3")
        cfg, err := LoadConfig(original)
        if err != nil {
                t.Fatalf("Load failed: %v", err)
        }
        marshaled, err := SaveConfig(cfg)
        if err != nil {
                t.Fatalf("Save failed: %v", err)
        }
        cfg2, err := LoadConfig(marshaled)
        if err != nil {
                t.Fatalf("Re-load failed: %v", err)
        }
        if len(cfg2) != 3 {
                t.Errorf("Round-trip length mismatch")
        }
        for i, item := range cfg2 {
                if item.Key != cfg[i].Key {
                        t.Errorf("Round-trip key order changed at position %d", i)
                }
        }
}

func TestMergeWithEmptySlices(t *testing.T) {
        base := MapSlice{{Key: "a", Value: 1}}
        emptySlice := MapSlice{}
        merged := Merge(base, emptySlice)
        if len(merged) != 1 {
                t.Errorf("Merge with empty override failed")
        }
        merged = Merge(emptySlice, base)
        if len(merged) != 1 {
                t.Errorf("Merge with empty base failed")
        }
}
