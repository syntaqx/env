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
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v %v", err, msgAndArgs)
	}
}

func assertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error, got nil %v", msgAndArgs...)
	}
}

func assertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v %v", expected, actual, msgAndArgs)
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
	setEnvForTest("DEBUG", "true")
	setEnvForTest("PORT", "9090")
	setEnvForTest("REDIS_HOST", "host1,host2")
	setEnvForTest("REDIS_MODE", "cluster")
	setEnvForTest("DATABASE_HOST", "dbhost")
	setEnvForTest("DATABASE_PORT", "5432")
	setEnvForTest("DATABASE_USERNAME", "admin")
	setEnvForTest("DATABASE_PASSWORD", "secret")
	setEnvForTest("DATABASE_NAME", "mydb")
	setEnvForTest("NESTED_FIELD", "nested_value")

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

func TestUnmarshalFloat(t *testing.T) {
	setEnvForTest("FLOAT32_VALUE", "3.14")
	setEnvForTest("FLOAT64_VALUE", "3.14159")

	var cfg struct {
		Float32Value float32 `env:"FLOAT32_VALUE"`
		Float64Value float64 `env:"FLOAT64_VALUE"`
	}
	err := Unmarshal(&cfg)
	assertNoError(t, err, "Unmarshal FloatConfig")

	expected := struct {
		Float32Value float32
		Float64Value float64
	}{
		Float32Value: 3.14,
		Float64Value: 3.14159,
	}

	assertEqual(t, expected.Float32Value, cfg.Float32Value, "Float32Value")
	assertEqual(t, expected.Float64Value, cfg.Float64Value, "Float64Value")

	err = Unset("FLOAT32_VALUE")
	assertNoError(t, err, "Unset FLOAT32_VALUE")

	err = Unset("FLOAT64_VALUE")
	assertNoError(t, err, "Unset FLOAT64_VALUE")
}

func TestUnmarshalUnsupportedKind(t *testing.T) {
	setEnvForTest("UNSUPPORTED", "invalid")

	var cfg struct {
		Unsupported complex64 `env:"UNSUPPORTED"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal Unsupported kind")

	err = Unset("UNSUPPORTED")
	assertNoError(t, err, "Unset UNSUPPORTED")
}

func TestUnmarshalSetFieldIntError(t *testing.T) {
	setEnvForTest("INVALID_INT", "invalid")

	var cfg struct {
		InvalidInt int `env:"INVALID_INT"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal InvalidInt")

	err = Unset("INVALID_INT")
	assertNoError(t, err, "Unset INVALID_INT")
}

func TestUnmarshalSetFieldFloatError(t *testing.T) {
	setEnvForTest("INVALID_FLOAT", "invalid")

	var cfg struct {
		InvalidFloat float64 `env:"INVALID_FLOAT"`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal InvalidFloat")

	err = Unset("INVALID_FLOAT")
	assertNoError(t, err, "Unset INVALID_FLOAT")
}

func TestUnmarshalSetFieldNestedError(t *testing.T) {
	type NestedConfig struct {
		NestedField string `env:"NESTED_FIELD,required"`
	}

	setEnvForTest("NESTED_FIELD", "") // Setting an empty value to trigger required error

	var cfg struct {
		Nested NestedConfig `env:""`
	}
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal NestedConfig")

	err = Unset("NESTED_FIELD")
	assertNoError(t, err, "Unset NESTED_FIELD")
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

	setEnvForTest("REQUIRED_VAR", "value")
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

	setEnvForTest("INVALID_UINT", "invalid")
	err := Unmarshal(&cfg)
	assertError(t, err, "Unmarshal InvalidUint")

	setEnvForTest("INVALID_FLOAT", "invalid")
	err = Unmarshal(&cfg)
	assertError(t, err, "Unmarshal InvalidFloat")

	setEnvForTest("UNSUPPORTED", "invalid")
	err = Unmarshal(&cfg)
	assertError(t, err, "Unmarshal Unsupported")
}

func TestSetFieldUint(t *testing.T) {
	type Config struct {
		UIntField uint `env:"UINT_FIELD"`
	}

	var cfg Config

	err := Set("UINT_FIELD", "42")
	assertNoError(t, err, "Set UINT_FIELD")

	field := reflect.ValueOf(&cfg).Elem().FieldByName("UIntField")
	err = setField(field, "42")
	assertNoError(t, err, "setField Uint")

	assertEqual(t, uint(42), cfg.UIntField, "UintField value")

	err = Unset("UINT_FIELD")
	assertNoError(t, err, "Unset UINT_FIELD")
}

func TestSetError(t *testing.T) {
	err := Set("", "value")
	assertError(t, err, "Set")
}

func setEnvForTest(t *testing.T, name string, value string) {
	err := Set(name, value)
	if err != nil {
		 t.Fatalf("setting env variable %s with %q: %v", name, value, err
	}

	t.Cleanup(func () {
		 err = Unset(name) // clean after the test
		 if err  != nil {
			  t.Fatalf("unsetting env variable %s: %v", name, err
		  }
	})
}
