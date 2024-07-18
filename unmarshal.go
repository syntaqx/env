package env

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Unmarshal reads environment variables into a struct based on `env` tags.
func Unmarshal(data interface{}) error {
	return unmarshalWithPrefix(data, "")
}

// unmarshalWithPrefix unmarshals environment variables into a struct with a given prefix.
func unmarshalWithPrefix(data interface{}, prefix string) error {
	v := reflect.ValueOf(data).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("env")

		// Handle nested structs with optional prefixes
		if field.Kind() == reflect.Struct {
			if err := unmarshalStruct(field.Addr().Interface(), prefix, tag); err != nil {
				return err
			}
			continue
		}

		if tag == "" {
			continue
		}

		if err := unmarshalField(field, prefix, tag); err != nil {
			return err
		}
	}

	return nil
}

// unmarshalStruct handles unmarshaling nested structs
func unmarshalStruct(data interface{}, prefix, tag string) error {
	newPrefix := prefix
	if tag != "" {
		newPrefix = prefix + tag + "_"
	}
	return unmarshalWithPrefix(data, newPrefix)
}

// unmarshalField handles unmarshaling individual fields based on tags
func unmarshalField(field reflect.Value, prefix, tag string) error {
	tagOpts := parseTag(tag)
	value, found := findFieldValue(tagOpts.keys, prefix)

	if !found {
		value = tagOpts.fallback
	}

	if tagOpts.required && value == "" {
		return fmt.Errorf("required environment variable %s is not set", tagOpts.keys[0])
	}

	return setField(field, value)
}

// findFieldValue tries to find environment variable value based on keys
func findFieldValue(keys []string, prefix string) (string, bool) {
	for _, key := range keys {
		fullKey := prefix + key
		if val, ok := Lookup(fullKey); ok {
			return val, true
		}
	}
	return "", false
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
		inBrackets := false
		start := 0
		for i := 0; i < len(extraParts); i++ {
			switch extraParts[i] {
			case '[':
				inBrackets = true
			case ']':
				inBrackets = false
			case ',':
				if !inBrackets {
					part := extraParts[start:i]
					start = i + 1
					parsePart(part, &fallbackValue, &required)
				}
			}
		}
		part := extraParts[start:]
		parsePart(part, &fallbackValue, &required)
	}

	return tagOptions{
		keys:     keys,
		fallback: fallbackValue,
		required: required,
	}
}

func parsePart(part string, fallbackValue *string, required *bool) {
	if strings.Contains(part, "default=[") || strings.Contains(part, "fallback=[") {
		re := regexp.MustCompile(`(?:default|fallback)=\[(.*?)]`)
		matches := re.FindStringSubmatch(part)
		if len(matches) > 1 {
			*fallbackValue = matches[1]
		}
	} else if strings.Contains(part, "default=") || strings.Contains(part, "fallback=") {
		re := regexp.MustCompile(`(?:default|fallback)=([^,]+)`)
		matches := re.FindStringSubmatch(part)
		if len(matches) > 1 {
			*fallbackValue = matches[1]
		}
	} else if strings.TrimSpace(part) == "required" {
		*required = true
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
