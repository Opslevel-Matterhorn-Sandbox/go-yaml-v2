# Configlayer

Configlayer loads YAML config and layers overrides on top of defaults. Key order is preserved when reading and writing so merged config stays readable and stable.

## Functionality

- **Layered config:** You provide a base (default) YAML config and an overrides YAML. Configlayer merges them so override values win for the same key, and keys only in overrides are appended. Base key order is kept; new keys from overrides keep their order. Useful for app config (defaults in repo) plus env-specific or user overrides (e.g. `env: production`, `replicas: 5`).

- **Load / Save:** `LoadConfig(data []byte)` unmarshals YAML into an order-preserving structure. `SaveConfig(cfg)` marshals it back to YAML. Round-tripping keeps key order.

- **Lookup:** `Get(cfg, key)` returns the value for a key and whether it was present.

- **Merge:** `Merge(base, overrides)` returns a new config: every key from base (with value from overrides when the key exists there), then every key that appears only in overrides, in order. Override values take precedence; key order is deterministic and stable for diffing and human editing.

The binary demonstrates the flow: it loads inline default and override YAML, merges them, and prints the result (e.g. `env` and `log_level` overridden, `replicas` added).

## Prerequisites

- Go 1.21+

### Installing Go

- **Official installer:** Download the latest Go release from [go.dev/download](https://go.dev/dl/) and run the installer for your OS.
- **macOS (Homebrew):** `brew install go`
- **Linux (apt):** `sudo apt update && sudo apt install golang-go`
- **Linux (snap):** `sudo snap install go --classic`

Verify the install: `go version` (should report 1.21 or higher).

## Run

```bash
go run .
```

Prints the merged config (defaults with overrides applied).

## Tests

Run tests for the current package (required so `config.go` is compiled with the tests):

```bash
go test .
```

Or from the repo root: `go test ./...`. Do not run `go test config_test.go` by itself—that compiles only the test file and will fail with undefined symbols.
