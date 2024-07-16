package env

import (
	"reflect"
	"runtime"
	"testing"
)

type RedisMode string

const (
	RedisModeStandalone RedisMode = "standalone"
	RedisModeCluster    RedisMode = "cluster"
)

type NestedConfig struct {
	NestedField string `env:"NESTED_FIELD,default=nested"`
}

type DatabaseConfig struct {
	Host     string       `env:"DATABASE_HOST,default=localhost"`
	Port     int          `env:"DATABASE_PORT|DB_PORT,fallback=3306"`
	Username string       `env:"DATABASE_USERNAME,default=root"`
	Password string       `env:"DATABASE_PASSWORD,required"`
	Database string       `env:"DATABASE_NAME"`
	Nested   NestedConfig `env:""`
}

type Config struct {
	Debug     bool           `env:"DEBUG"`
	Port      string         `env:"PORT,default=8080"`
	RedisHost []string       `env:"REDIS_HOST|REDIS_HOSTS,default=localhost:6379"`
	RedisMode RedisMode      `env:"REDIS_MODE,default=standalone"`
	Database  DatabaseConfig `env:""`
}

func assertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func assertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func assertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestSetUnset(t *testing.T) {
	key, value := "TEST_KEY", "TEST_VALUE"
	err := Set(key, value)
	assertNoError(t, err, "Set")

	actual := Get(key)
	assertEqual(t, value, actual, "Get")

	err = Unset(key)
	assertNoError(t, err, "Unset")

	if _, ok := Lookup(key); ok {
		t.Errorf("Lookup: expected %s to be unset", key)
	}
}

func TestEnvFunctions(t *testing.T) {
	tests := []struct {
		name                string
		key                 string
		setValue            string
		fallback            interface{}
		expected            interface{}
		expectedFallback    interface{}
		setFunc             func(string, string) error
		getFunc             func(string) (interface{}, error)
		getFuncWithFallback func(string, interface{}) interface{}
	}{
		{
			name:                "string",
			key:                 "TEST_STRING",
			setValue:            "stringValue",
			fallback:            "fallbackValue",
			expected:            "stringValue",
			expectedFallback:    "fallbackValue",
			setFunc:             Set,
			getFunc:             func(key string) (interface{}, error) { return Get(key), nil },
			getFuncWithFallback: func(key string, fallback interface{}) interface{} { return GetWithFallback(key, fallback.(string)) },
		},
		{
			name:                "int",
			key:                 "TEST_INT",
			setValue:            "42",
			fallback:            10,
			expected:            42,
			expectedFallback:    10,
			setFunc:             Set,
			getFunc:             func(key string) (interface{}, error) { return GetInt(key) },
			getFuncWithFallback: func(key string, fallback interface{}) interface{} { return GetIntWithFallback(key, fallback.(int)) },
		},
		{
			name:                "bool",
			key:                 "TEST_BOOL",
			setValue:            "true",
			fallback:            false,
			expected:            true,
			expectedFallback:    false,
			setFunc:             Set,
			getFunc:             func(key string) (interface{}, error) { return GetBool(key), nil },
			getFuncWithFallback: func(key string, fallback interface{}) interface{} { return GetBoolWithFallback(key, fallback.(bool)) },
		},
		{
			name:             "float",
			key:              "TEST_FLOAT",
			setValue:         "42.42",
			fallback:         10.1,
			expected:         42.42,
			expectedFallback: 10.1,
			setFunc:          Set,
			getFunc:          func(key string) (interface{}, error) { return GetFloat(key) },
			getFuncWithFallback: func(key string, fallback interface{}) interface{} {
				return GetFloatWithFallback(key, fallback.(float64))
			},
		},
		{
			name:             "slice",
			key:              "TEST_SLICE",
			setValue:         "value1,value2",
			fallback:         []string{"fallback1", "fallback2"},
			expected:         []string{"value1", "value2"},
			expectedFallback: []string{"fallback1", "fallback2"},
			setFunc:          Set,
			getFunc:          func(key string) (interface{}, error) { return GetSlice(key) },
			getFuncWithFallback: func(key string, fallback interface{}) interface{} {
				return GetSliceWithFallback(key, fallback.([]string))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with set value
			err := tt.setFunc(tt.key, tt.setValue)
			assertNoError(t, err, tt.name)

			val, err := tt.getFunc(tt.key)
			assertNoError(t, err, tt.name)
			assertEqual(t, tt.expected, val, tt.name)

			val = tt.getFuncWithFallback(tt.key, tt.fallback)
			assertEqual(t, tt.expected, val, tt.name)

			err = Unset(tt.key)
			assertNoError(t, err, tt.name)

			// Test with fallback value
			val = tt.getFuncWithFallback(tt.key, tt.fallback)
			assertEqual(t, tt.expectedFallback, val, tt.name)
		})
	}
}

func TestRequire(t *testing.T) {
	key := "TEST_REQUIRED"

	err := Require(key)
	assertError(t, err, "Require")

	err = Set(key, "value")
	assertNoError(t, err, "Set")

	err = Require(key)
	assertNoError(t, err, "Require")

	err = Unset(key)
	assertNoError(t, err, "Unset")
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"1", true},
		{"yes", true},
		{"false", false},
		{"0", false},
		{"no", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		result := parseBool(tt.input)
		assertEqual(t, tt.expected, result, "parseBool")
	}
}

func TestUnsetError(t *testing.T) {
	err := Unset("")
	if runtime.GOOS == "windows" {
		assertError(t, err, "Unset")
	} else {
		assertNoError(t, err, "Unset")
	}
}

func TestUnmarshal(t *testing.T) {
	_ = Set("DEBUG", "true")
	_ = Set("PORT", "9090")
	_ = Set("REDIS_HOST", "host1,host2")
	_ = Set("REDIS_MODE", "cluster")
	_ = Set("DATABASE_HOST", "dbhost")
	_ = Set("DATABASE_PORT", "5432")
	_ = Set("DATABASE_USERNAME", "admin")
	_ = Set("DATABASE_PASSWORD", "secret")
	_ = Set("DATABASE_NAME", "mydb")
	_ = Set("NESTED_FIELD", "nested_value")

	var cfg Config
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal")

	expected := Config{
		Debug: true,
		Port:  "9090",
		RedisHost: []string{
			"host1",
			"host2",
		},
		RedisMode: RedisModeCluster,
		Database: DatabaseConfig{
			Host:     "dbhost",
			Port:     5432,
			Username: "admin",
			Password: "secret",
			Database: "mydb",
			Nested: NestedConfig{
				NestedField: "nested_value",
			},
		},
	}

	assertEqual(t, expected, cfg, "Unmarshal")
}

func TestUnmarshalRequired(t *testing.T) {
	type RequiredConfig struct {
		RequiredVar string `env:"REQUIRED_VAR,required"`
	}

	var cfg RequiredConfig
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal Required")

	expectedErr := "required environment variable REQUIRED_VAR is not set"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("expected error %s, got %s", expectedErr, err.Error())
	}

	_ = Set("REQUIRED_VAR", "value")
	err = Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal Required")
	assertEqual(t, "value", cfg.RequiredVar, "Unmarshal Required")
}

func TestUnmarshalSetFieldErrors(t *testing.T) {
	type InvalidConfig struct {
		InvalidUint  uint      `env:"INVALID_UINT"`
		InvalidFloat float32   `env:"INVALID_FLOAT"`
		Unsupported  complex64 `env:"UNSUPPORTED"`
	}

	var cfg InvalidConfig

	_ = Set("INVALID_UINT", "invalid")
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal InvalidUint")

	_ = Set("INVALID_FLOAT", "invalid")
	err = Unmarshal(&cfg)
	assertError(t, err, "Unmarshal InvalidFloat")

	_ = Set("UNSUPPORTED", "invalid")
	err = Unmarshal(&cfg)
	assertError(t, err, "Unmarshal Unsupported")
}

func TestSetError(t *testing.T) {
	err := Set("", "value")
	assertError(t, err, "Set")
}
