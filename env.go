package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
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

// GetSlice returns the value of a comma-separated environment variable as a slice of strings.
func GetSlice(key string) ([]string, error) {
	value := Get(key)
	if value == "" {
		return nil, fmt.Errorf("environment variable %s not set", key)
	}
	return strings.Split(value, ","), nil
}

// GetSliceWithFallback returns the value of a comma-separated environment variable as a slice
// of strings or a fallback value if the environment variable is not set.
func GetSliceWithFallback(key string, fallback []string) []string {
	value, err := GetSlice(key)
	if err != nil {
		return fallback
	}
	return value
}

// GetBool returns the value of an environment variable as a boolean.
func GetBool(key string) bool {
	return parseBool(Get(key))
}

// GetBoolWithFallback returns the value of an environment variable as a boolean
// or a fallback value if the environment variable is not set.
func GetBoolWithFallback(key string, fallback bool) bool {
	if value, ok := Lookup(key); ok {
		return parseBool(value)
	}
	return fallback
}

// GetInt returns the value of an environment variable as an integer.
func GetInt(key string) (int, error) {
	value := Get(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s not set", key)
	}
	return strconv.Atoi(value)
}

// GetIntWithFallback returns the value of an environment variable as an integer
// or a fallback value if the environment variable is not set or invalid.
func GetIntWithFallback(key string, fallback int) int {
	value, err := GetInt(key)
	if err != nil {
		return fallback
	}
	return value
}

// GetFloat returns the value of an environment variable as a float.
func GetFloat(key string) (float64, error) {
	value := Get(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s not set", key)
	}
	return strconv.ParseFloat(value, 64)
}

// GetFloatWithFallback returns the value of an environment variable as a float
// or a fallback value if the environment variable is not set or invalid.
func GetFloatWithFallback(key string, fallback float64) float64 {
	value, err := GetFloat(key)
	if err != nil {
		return fallback
	}
	return value
}

// Require checks if an environment variable is set and returns an error if it is not.
func Require(key string) error {
	if _, ok := Lookup(key); !ok {
		return fmt.Errorf("required environment variable %s is not set", key)
	}
	return nil
}

// parseBool is a helper function to parse a boolean from a string value.
func parseBool(value string) bool {
	switch strings.ToLower(value) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return false
	}
}

// Unmarshal reads environment variables into a struct based on `env` tags.
func Unmarshal(data any) error {
	return unmarshalWithPrefix(data, "")
}

// unmarshalWithPrefix unmarshals environment variables into a struct with a given prefix.
func unmarshalWithPrefix(data any, prefix string) error {
	v := reflect.ValueOf(data).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("env")

		// Handle nested structs with optional prefixes
		if field.Kind() == reflect.Struct {
			newPrefix := prefix
			if tag != "" {
				newPrefix = prefix + tag + "_"
			}
			if err := unmarshalWithPrefix(field.Addr().Interface(), newPrefix); err != nil {
				return err
			}
			continue
		}

		if tag == "" {
			continue
		}

		tagOpts := parseTag(tag)
		var value string
		var found bool
		for _, key := range tagOpts.keys {
			fullKey := prefix + key
			if val, ok := Lookup(fullKey); ok {
				value = val
				found = true
				break
			}
		}

		if !found {
			value = tagOpts.fallback
		}

		if tagOpts.required && value == "" {
			return fmt.Errorf("required environment variable %s is not set", tagOpts.keys[0])
		}

		if err := setField(field, value); err != nil {
			return err
		}
	}

	return nil
}

// tagOptions holds parsed tag options
type tagOptions struct {
	keys     []string
	fallback string
	required bool
}

// parseTag parses the struct tag into tagOptions
func parseTag(tag string) tagOptions {
	parts := strings.Split(tag, ",")
	keys := strings.Split(parts[0], "|")
	var fallbackValue string
	required := false
	if len(parts) > 1 {
		for _, part := range parts[1:] {
			if strings.HasPrefix(part, "default=") {
				fallbackValue = strings.TrimPrefix(part, "default=")
			}
			if strings.HasPrefix(part, "fallback=") {
				fallbackValue = strings.TrimPrefix(part, "fallback=")
			}
			if part == "required" {
				required = true
			}
		}
	}

	return tagOptions{
		keys:     keys,
		fallback: fallbackValue,
		required: required,
	}
}

// setField sets the value of a struct field based on its type
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		field.SetBool(parseBool(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(strings.Split(value, ",")))
		}
	default:
		return fmt.Errorf("unsupported kind %s", field.Kind())
	}
	return nil
}
