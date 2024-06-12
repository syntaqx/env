package env

import (
	"os"

	"github.com/joho/godotenv"
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

// Load will read your env file(s) and load them into ENV for this process.
//
// Call this function as close as possible to the start of your program
// (ideally in main).
//
// If you call Load without any args it will default to loading .env in the
// current path.
//
// You can otherwise tell it which files to load (there can be more than one)
// like:
//
//	env.Load("fileone", "filetwo")
//
// > [!IMPORTANT] This __WILL NOT__ override an env variable that already
// > exists. Consider the .env file to set dev vars or sensible defaults.
func Load(filenames ...string) (err error) {
	return godotenv.Load(filenames...)
}

// Overload will read your env file(s) and load them into ENV for this process.
//
// Call this function as close as possible to the start of your program
// (ideally in main).
//
// If you call Load without any args it will default to loading .env in the
// current path.
//
// You can otherwise tell it which files to load (there can be more than one)
// like:
//
//	env.Overload("fileone", "filetwo")
//
// > [!IMPORTANT] This __WILL__ override an env variable that already
// > exists. Consider the .env file to forcefully set all vars.
func Overload(filenames ...string) (err error) {
	return godotenv.Overload(filenames...)
}
