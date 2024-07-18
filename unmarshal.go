package env

import (
	"fmt"
	"reflect"
	"regexp"
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
	parts := strings.SplitN(tag, ",", 2)
	keys := strings.Split(parts[0], "|")
	var fallbackValue string
	required := false
	if len(parts) > 1 {
		extraParts := parts[1]
		if strings.Contains(extraParts, "default=[") || strings.Contains(extraParts, "fallback=[") {
			re := regexp.MustCompile(`(?:default|fallback)=\[(.*?)\]`)
			matches := re.FindStringSubmatch(extraParts)
			if len(matches) > 1 {
				fallbackValue = matches[1]
			}
		} else if strings.Contains(extraParts, "default=") || strings.Contains(extraParts, "fallback=") {
			re := regexp.MustCompile(`(?:default|fallback)=([^,]+)`)
			matches := re.FindStringSubmatch(extraParts)
			if len(matches) > 1 {
				fallbackValue = matches[1]
			}
		}
		required = strings.Contains(extraParts, "required")
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
