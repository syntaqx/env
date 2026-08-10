// Package env provides utilities for reading, writing, and parsing environment
// variables.
//
// It offers typed getters (string, bool, int, uint, float, time.Duration, and
// their slice variants) with optional fallbacks, plus Unmarshal for loading
// environment variables into structs via `env` tags. Unmarshal supports nested
// struct prefixes, default and required values, variable expansion, loading
// values from files, and any type implementing encoding.TextUnmarshaler.
package env

import (
	"fmt"
	"os"
)

// Set sets an environment variable.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset unsets an environment variable.
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Lookup returns the value of an environment variable and a boolean indicating
// whether the variable is present in the environment.
func Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Require checks if an environment variable is set and returns an error if it is not.
func Require(key string) error {
	if _, ok := Lookup(key); !ok {
		return fmt.Errorf("required environment variable %s is not set", key)
	}
	return nil
}
