package env

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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
// setField sets the value of a struct field based on its type
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		boolValue, err := parseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
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
		elemType := field.Type().Elem()
		switch elemType.Kind() {
		case reflect.String:
			field.Set(reflect.ValueOf(strings.Split(value, ",")))
		case reflect.Bool:
			boolSlice, err := parseBoolSlice(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(boolSlice))
		case reflect.Int:
			intSlice, err := parseIntSlice(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(intSlice))
		case reflect.Uint:
			uintSlice, err := parseUintSlice(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(uintSlice))
		case reflect.Float64:
			floatSlice, err := parseFloatSlice(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(floatSlice))
		default:
			return fmt.Errorf("unsupported slice element kind %s", elemType.Kind())
		}
	default:
		return fmt.Errorf("unsupported kind %s", field.Kind())
	}
	return nil
}

// parseBoolSlice parses a comma-separated string into a slice of bools
func parseBoolSlice(value string) ([]bool, error) {
	values := strings.Split(value, ",")
	result := make([]bool, len(values))
	for i, v := range values {
		boolValue, err := parseBool(v)
		if err != nil {
			return nil, err
		}
		result[i] = boolValue
	}
	return result, nil
}

// parseIntSlice parses a comma-separated string into a slice of ints
func parseIntSlice(value string) ([]int, error) {
	values := strings.Split(value, ",")
	result := make([]int, len(values))
	for i, v := range values {
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		result[i] = intValue
	}
	return result, nil
}

// parseUintSlice parses a comma-separated string into a slice of uints
func parseUintSlice(value string) ([]uint, error) {
	values := strings.Split(value, ",")
	result := make([]uint, len(values))
	for i, v := range values {
		uintValue, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		result[i] = uint(uintValue)
	}
	return result, nil
}

// parseFloatSlice parses a comma-separated string into a slice of floats
func parseFloatSlice(value string) ([]float64, error) {
	values := strings.Split(value, ",")
	result := make([]float64, len(values))
	for i, v := range values {
		floatValue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		result[i] = floatValue
	}
	return result, nil
}
