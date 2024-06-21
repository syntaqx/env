package env

import (
	"os"
	"testing"
)

func TestGetWithFallback(t *testing.T) {
	tests := []struct {
		setValue  bool
		key       string
		fallback  string
		envValue  string
		wantValue string
	}{
		{
			setValue:  true,
			key:       "EXISTING_KEY",
			fallback:  "fallback_value",
			envValue:  "existing_value",
			wantValue: "existing_value",
		},
		{
			setValue:  false,
			key:       "NON_EXISTING_KEY",
			fallback:  "fallback_value",
			envValue:  "",
			wantValue: "fallback_value",
		},
		{
			setValue:  true,
			key:       "EMPTY_VALUE_KEY",
			fallback:  "fallback_value",
			envValue:  "",
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			// Set the environment variable
			if tt.setValue {
				os.Setenv(tt.key, tt.envValue)
			}

			// Call the function under test
			got := GetWithFallback(tt.key, tt.fallback)

			// Check if the returned value matches the expected value
			if got != tt.wantValue {
				t.Errorf("GetWithFallback(%q, %q) = %q, want %q", tt.key, tt.fallback, got, tt.wantValue)
			}

			// Unset the environment variable
			os.Unsetenv(tt.key)
		})
	}
}
func TestGetBoolWithFallback(t *testing.T) {
	tests := []struct {
		setValue  bool
		key       string
		fallback  bool
		envValue  string
		wantValue bool
	}{
		{
			setValue:  true,
			key:       "EXISTING_KEY",
			fallback:  true,
			envValue:  "true",
			wantValue: true,
		},
		{
			setValue:  false,
			key:       "NON_EXISTING_KEY",
			fallback:  true,
			envValue:  "",
			wantValue: true,
		},
		{
			setValue:  true,
			key:       "EMPTY_VALUE_KEY",
			fallback:  false,
			envValue:  "",
			wantValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			// Set the environment variable
			if tt.setValue {
				os.Setenv(tt.key, tt.envValue)
			}
			// Call the function under test
			got := GetBoolWithFallback(tt.key, tt.fallback)
			// Check if the returned value matches the expected value
			if got != tt.wantValue {
				t.Errorf("GetBoolWithFallback(%q, %v) = %v, want %v", tt.key, tt.fallback, got, tt.wantValue)
			}
			// Unset the environment variable
			os.Unsetenv(tt.key)
		})
	}
}
