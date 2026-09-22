package config

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

// Config holds all service configuration.
type Config struct {
	Addr   string
	Secret string
}

// Load reads configuration from defaults, env vars, and flags.
// Precedence: defaults < env vars < flags.
func Load() *Config {
	c := &Config{
		Addr:   ":8787",
		Secret: "",
	}

	// Env vars
	if v := os.Getenv("RANDOMKIT_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("RANDOMKIT_SECRET"); v != "" {
		c.Secret = v
	}

	// Flags
	flag.StringVar(&c.Addr, "addr", c.Addr, "listen address")
	flag.StringVar(&c.Secret, "secret", c.Secret, "token signing secret (auto-generated if empty)")
	flag.Parse()

	// Auto-generate secret if not provided
	if c.Secret == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			panic(fmt.Sprintf("failed to generate secret: %v", err))
		}
		c.Secret = hex.EncodeToString(b)
	}

	return c
}
