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

// GetSliceWithFallback returns the value of a comma-separated environment variable as a slice
// of strings or a fallback value if the environment variable is not set.
func GetSliceWithFallback(key string, fallback []string) []string {
	if value, ok := Lookup(key); ok {
		return strings.Split(value, ",")
	}
	return fallback
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
	value, err := strconv.Atoi(Get(key))
	if err != nil {
		return 0, fmt.Errorf("error converting %s to integer: %w", key, err)
	}
	return value, nil
}

// GetIntWithFallback returns the value of an environment variable as an integer
// or a fallback value if the environment variable is not set or invalid.
func GetIntWithFallback(key string, fallback int) int {
	if value, ok := Lookup(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// GetFloat returns the value of an environment variable as a float.
func GetFloat(key string) (float64, error) {
	value, err := strconv.ParseFloat(Get(key), 64)
	if err != nil {
		return 0, fmt.Errorf("error converting %s to float: %w", key, err)
	}
	return value, nil
}

// GetFloatWithFallback returns the value of an environment variable as a float
// or a fallback value if the environment variable is not set or invalid.
func GetFloatWithFallback(key string, fallback float64) float64 {
	if value, ok := Lookup(key); ok {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return fallback
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
func Unmarshal(cfg interface{}) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("env")

		if tag == "" {
			if field.Kind() == reflect.Struct {
				if err := Unmarshal(field.Addr().Interface()); err != nil {
					return err
				}
			}
			continue
		}

		parts := strings.Split(tag, ",")
		keys := strings.Split(parts[0], "|")
		var defaultValue string
		if len(parts) > 1 && strings.HasPrefix(parts[1], "default=") {
			defaultValue = strings.TrimPrefix(parts[1], "default=")
		}

		var value string
		var found bool
		for _, key := range keys {
			if val, ok := Lookup(key); ok {
				value = val
				found = true
				break
			}
		}

		if !found {
			value = defaultValue
		}

		if err := setField(field, value); err != nil {
			return err
		}
	}

	return nil
}

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
