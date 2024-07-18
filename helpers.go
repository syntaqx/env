package env

import (
	"fmt"
	"strconv"
	"strings"
)

func parseBool(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value %s", value)
	}
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
