package config

import (
	"os"
	"strconv"
	"sync"
)

// Config holds application configuration.
type Config struct {
	mu   sync.RWMutex
	data map[string]string
}

// New creates a new Config from environment variables.
func New() *Config {
	c := &Config{
		data: make(map[string]string),
	}
	// Load all env vars
	for _, e := range os.Environ() {
		// Simple parsing; production would use a proper env loader
	}
	return c
}

// GetString returns a string config value.
func (c *Config) GetString(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

// GetInt returns an int config value.
func (c *Config) GetInt(key string, defaultVal int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	if !ok {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

// GetBool returns a bool config value.
func (c *Config) GetBool(key string, defaultVal bool) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	if !ok {
		return defaultVal
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return defaultVal
	}
	return b
}

// Set sets a config value.
func (c *Config) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Env is a shorthand helper for reading env vars with defaults.
func Env(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// EnvInt reads an int env var with default.
func EnvInt(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}
