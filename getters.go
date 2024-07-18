package env

import (
	"fmt"
	"strconv"
	"strings"
)

// Get returns the value of an environment variable.
func Get(key string) (string, error) {
	if value, ok := Lookup(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s not set", key)
}

// GetWithFallback returns the value of an environment variable or a fallback
// value if the environment variable is not set.
func GetWithFallback(key string, fallback string) string {
	if value, err := Get(key); err == nil {
		return value
	}
	return fallback
}

// GetBool returns the value of an environment variable as a boolean.
func GetBool(key string) (bool, error) {
	value, err := Get(key)
	if err != nil {
		return false, err
	}
	return parseBool(value)
}

// GetBoolWithFallback returns the value of an environment variable as a boolean
// or a fallback value if the environment variable is not set.
func GetBoolWithFallback(key string, fallback bool) (bool, error) {
	if value, err := GetBool(key); err == nil {
		return value, nil
	}
	return fallback, nil
}

// GetInt returns the value of an environment variable as an integer.
func GetInt(key string) (int, error) {
	value, err := Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

// GetIntWithFallback returns the value of an environment variable as an integer
// or a fallback value if the environment variable is not set or invalid.
func GetIntWithFallback(key string, fallback int) (int, error) {
	if value, err := GetInt(key); err == nil {
		return value, nil
	}
	return fallback, nil
}

// GetFloat returns the value of an environment variable as a float.
func GetFloat(key string) (float64, error) {
	value, err := Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(value, 64)
}

// GetFloatWithFallback returns the value of an environment variable as a float
// or a fallback value if the environment variable is not set or invalid.
func GetFloatWithFallback(key string, fallback float64) (float64, error) {
	if value, err := GetFloat(key); err == nil {
		return value, nil
	}
	return fallback, nil
}

// GetStringSlice returns the value of a comma-separated environment variable as a slice of strings.
func GetStringSlice(key string) ([]string, error) {
	value, err := Get(key)
	if err != nil {
		return nil, err
	}
	return strings.Split(value, ","), nil
}

// GetStringSliceWithFallback returns the value of a comma-separated environment variable as a slice
// of strings or a fallback value if the environment variable is not set.
func GetStringSliceWithFallback(key string, fallback []string) ([]string, error) {
	if value, err := GetStringSlice(key); err == nil {
		return value, nil
	}
	return fallback, nil
}
