package main

import (
	"fmt"
	"os"
)

func main() {
	defaults := `name: myapp
version: "1.0"
env: dev
log_level: info`
	overrides := `env: production
log_level: warn
replicas: 3`

	base, err := LoadConfig([]byte(defaults))
	if err != nil {
		fmt.Fprintln(os.Stderr, "load defaults:", err)
		os.Exit(1)
	}
	ov, err := LoadConfig([]byte(overrides))
	if err != nil {
		fmt.Fprintln(os.Stderr, "load overrides:", err)
		os.Exit(1)
	}
	merged := Merge(base, ov)
	out, err := SaveConfig(merged)
	if err != nil {
		fmt.Fprintln(os.Stderr, "save config:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
