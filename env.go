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
