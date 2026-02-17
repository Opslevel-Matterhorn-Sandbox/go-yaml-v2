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
        if len(cfg) != 3 {
                t.Errorf("expected 3 items, got %d", len(cfg))
        }
        if cfg[0].Key != "a" || cfg[1].Key != "b" || cfg[2].Key != "c" {
                t.Errorf("key order not preserved: %v", cfg)
        }
}

func TestSaveConfig(t *testing.T) {
        cfg := MapSlice{
                {Key: "name", Value: "test"},
                {Key: "version", Value: "2.0"},
        }
        out, err := SaveConfig(cfg)
        if err != nil {
                t.Fatalf("SaveConfig: %v", err)
        }
        roundTrip, err := LoadConfig(out)
        if err != nil {
                t.Fatalf("LoadConfig round-trip: %v", err)
        }
        if len(roundTrip) != 2 {
                t.Errorf("round-trip length: got %d", len(roundTrip))
        }
}

func TestGet(t *testing.T) {
        cfg := MapSlice{
                {Key: "name", Value: "myapp"},
                {Key: "count", Value: 42},
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
        base := MapSlice{
                {Key: "name", Value: "app"},
                {Key: "env", Value: "dev"},
                {Key: "replicas", Value: 1},
        }
        ov := MapSlice{
                {Key: "env", Value: "prod"},
                {Key: "replicas", Value: 5},
        }
        merged := Merge(base, ov)
        if len(merged) != 3 {
                t.Fatalf("merged length: got %d", len(merged))
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
        if merged[0].Key != "name" || merged[1].Key != "env" || merged[2].Key != "replicas" {
                t.Errorf("merge order: %v", merged)
        }
}

func TestMergeAddNewKeys(t *testing.T) {
        base := MapSlice{
                {Key: "a", Value: 1},
        }
        ov := MapSlice{
                {Key: "b", Value: 2},
                {Key: "c", Value: 3},
        }
        merged := Merge(base, ov)
        if len(merged) != 3 {
                t.Fatalf("merged length: got %d", len(merged))
        }
        if merged[0].Key != "a" || merged[1].Key != "b" || merged[2].Key != "c" {
                t.Errorf("merge order: %v", merged)
        }
}

func TestMapItemPreservesOrder(t *testing.T) {
        var slice MapSlice
        slice = append(slice, MapItem{Key: "first", Value: 1})
        slice = append(slice, MapItem{Key: "second", Value: 2})
        out, err := yaml.Marshal(slice)
        if err != nil {
                t.Fatal(err)
        }
        back, err := LoadConfig(out)
        if err != nil {
                t.Fatal(err)
        }
        if back[0].Key != "first" || back[1].Key != "second" {
                t.Errorf("order changed: %v", back)
        }
}
