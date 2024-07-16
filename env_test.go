package env

import (
	"reflect"
	"testing"
)

type RedisMode string

const (
	RedisModeStandalone RedisMode = "standalone"
	RedisModeCluster    RedisMode = "cluster"
)

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST,default=localhost"`
	Port     int    `env:"DATABASE_PORT|DB_PORT,fallback=3306"`
	Username string `env:"DATABASE_USERNAME,default=root"`
	Password string `env:"DATABASE_PASSWORD"`
	Database string `env:"DATABASE_NAME"`
}

type Config struct {
	Debug     bool           `env:"DEBUG"`
	Port      string         `env:"PORT,default=8080"`
	RedisHost []string       `env:"REDIS_HOST|REDIS_HOSTS,default=localhost:6379"`
	RedisMode RedisMode      `env:"REDIS_MODE,default=standalone"`
	Database  DatabaseConfig `env:""`
}

func TestSetUnset(t *testing.T) {
	key, value := "TEST_KEY", "TEST_VALUE"
	err := Set(key, value)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if Get(key) != value {
		t.Fatalf("expected %s, got %s", value, Get(key))
	}

	err = Unset(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := Lookup(key); ok {
		t.Fatalf("expected %s to be unset", key)
	}
}

func TestEnvFunctions(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		setValue   string
		fallback   interface{}
		expected   interface{}
		setFunc    func(string, string) error
		getFunc    func(string) (interface{}, error)
		unsetFunc  func(string) error
		lookupFunc func(string) (string, bool)
	}{
		{
			name:       "string",
			key:        "TEST_STRING",
			setValue:   "stringValue",
			fallback:   "fallbackValue",
			expected:   "stringValue",
			setFunc:    Set,
			getFunc:    func(key string) (interface{}, error) { return GetWithFallback(key, "fallbackValue"), nil },
			unsetFunc:  Unset,
			lookupFunc: Lookup,
		},
		{
			name:       "int",
			key:        "TEST_INT",
			setValue:   "42",
			fallback:   10,
			expected:   42,
			setFunc:    Set,
			getFunc:    func(key string) (interface{}, error) { return GetIntWithFallback(key, 10), nil },
			unsetFunc:  Unset,
			lookupFunc: Lookup,
		},
		{
			name:       "bool",
			key:        "TEST_BOOL",
			setValue:   "true",
			fallback:   false,
			expected:   true,
			setFunc:    Set,
			getFunc:    func(key string) (interface{}, error) { return GetBoolWithFallback(key, false), nil },
			unsetFunc:  Unset,
			lookupFunc: Lookup,
		},
		{
			name:       "float",
			key:        "TEST_FLOAT",
			setValue:   "42.42",
			fallback:   10.1,
			expected:   42.42,
			setFunc:    Set,
			getFunc:    func(key string) (interface{}, error) { return GetFloatWithFallback(key, 10.1), nil },
			unsetFunc:  Unset,
			lookupFunc: Lookup,
		},
		{
			name:     "slice",
			key:      "TEST_SLICE",
			setValue: "value1,value2",
			fallback: []string{"fallback1", "fallback2"},
			expected: []string{"value1", "value2"},
			setFunc:  Set,
			getFunc: func(key string) (interface{}, error) {
				return GetSliceWithFallback(key, []string{"fallback1", "fallback2"}), nil
			},
			unsetFunc:  Unset,
			lookupFunc: Lookup,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setFunc(tt.key, tt.setValue)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			val, err := tt.getFunc(tt.key)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(val, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, val)
			}

			err = tt.unsetFunc(tt.key)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if _, ok := tt.lookupFunc(tt.key); ok {
				t.Fatalf("expected %s to be unset", tt.key)
			}
		})
	}
}

func TestRequire(t *testing.T) {
	key := "TEST_REQUIRED"

	err := Require(key)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := Set(key, "value"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = Require(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := Unset(key); err != nil {
		t.Fatalf("expected no error, got %v", err)
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

	var cfg Config
	if err := Unmarshal(&cfg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

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
		},
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Fatalf("expected %+v, got %+v", expected, cfg)
	}
}

func TestMarshal(t *testing.T) {
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
		},
	}

	if err := Marshal(&expected); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var cfg Config
	if err := Unmarshal(&cfg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Fatalf("expected %+v, got %+v", expected, cfg)
	}
}
