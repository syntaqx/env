package env

import (
	"testing"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		err      bool
	}{
		{"true", true, false},
		{"1", true, false},
		{"yes", true, false},
		{"false", false, false},
		{"0", false, false},
		{"no", false, false},
		{"invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseBool(tt.input)
			if tt.err {
				assertError(t, err, "parseBool")
			} else {
				assertNoError(t, err, "parseBool")
				assertEqual(t, tt.expected, result, "parseBool")
			}
		})
	}
}
