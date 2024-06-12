package env

import (
	"os"
)

// Set sets an environment variable.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Lookup returns the value of an environment variable and a boolean indicating
// whether the variable is present in the environment.
func Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

// Get returns the value of an environment variable.
func Get(key string) string {
	return os.Getenv(key)
}

// GetWithFallback returns the value of an environment variable or a fallback
// value if the environment variable is not set.
func GetWithFallback(key string, fallback string) string {
	if value, ok := Lookup(key); ok {
		return value
	}
	return fallback
}
