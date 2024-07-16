package env_test

import (
	"reflect"
	"testing"

	"github.com/syntaqx/env"
)

func TestSetUnset(t *testing.T) {
	key, value := "TEST_KEY", "TEST_VALUE"
	err := env.Set(key, value)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if env.Get(key) != value {
		t.Fatalf("expected %s, got %s", value, env.Get(key))
	}

	err = env.Unset(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := env.Lookup(key); ok {
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
			setFunc:    env.Set,
			getFunc:    func(key string) (interface{}, error) { return env.GetWithFallback(key, "fallbackValue"), nil },
			unsetFunc:  env.Unset,
			lookupFunc: env.Lookup,
		},
		{
			name:       "int",
			key:        "TEST_INT",
			setValue:   "42",
			fallback:   10,
			expected:   42,
			setFunc:    env.Set,
			getFunc:    func(key string) (interface{}, error) { return env.GetIntWithFallback(key, 10), nil },
			unsetFunc:  env.Unset,
			lookupFunc: env.Lookup,
		},
		{
			name:       "bool",
			key:        "TEST_BOOL",
			setValue:   "true",
			fallback:   false,
			expected:   true,
			setFunc:    env.Set,
			getFunc:    func(key string) (interface{}, error) { return env.GetBoolWithFallback(key, false), nil },
			unsetFunc:  env.Unset,
			lookupFunc: env.Lookup,
		},
		{
			name:       "float",
			key:        "TEST_FLOAT",
			setValue:   "42.42",
			fallback:   10.1,
			expected:   42.42,
			setFunc:    env.Set,
			getFunc:    func(key string) (interface{}, error) { return env.GetFloatWithFallback(key, 10.1), nil },
			unsetFunc:  env.Unset,
			lookupFunc: env.Lookup,
		},
		{
			name:     "slice",
			key:      "TEST_SLICE",
			setValue: "value1,value2",
			fallback: []string{"fallback1", "fallback2"},
			expected: []string{"value1", "value2"},
			setFunc:  env.Set,
			getFunc: func(key string) (interface{}, error) {
				return env.GetSliceWithFallback(key, []string{"fallback1", "fallback2"}), nil
			},
			unsetFunc:  env.Unset,
			lookupFunc: env.Lookup,
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

	err := env.Require(key)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	env.Set(key, "value")
	err = env.Require(key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	env.Unset(key)
}

func TestUnmarshal(t *testing.T) {
	type RedisMode string

	const (
		RedisModeStandalone RedisMode = "standalone"
		RedisModeCluster    RedisMode = "cluster"
	)

	type DatabaseConfig struct {
		Host     string `env:"DATABASE_HOST|DB_HOST,default=localhost"`
		Port     int    `env:"DATABASE_PORT|DB_PORT,default=3306"`
		Username string `env:"DATABASE_USERNAME|DB_USER,default=root"`
		Password string `env:"DATABASE_PASSWORD|DB_PASS"`
		Database string `env:"DATABASE_NAME|DB_NAME"`
	}

	type Config struct {
		Debug     bool           `env:"DEBUG"`
		Port      string         `env:"PORT,default=8080"`
		RedisHost []string       `env:"REDIS_HOST|REDIS_HOSTS,default=localhost:6379"`
		RedisMode RedisMode      `env:"REDIS_MODE,default=standalone"`
		Database  DatabaseConfig `env:""`
	}

	env.Set("DEBUG", "true")
	env.Set("PORT", "9090")
	env.Set("REDIS_HOST", "host1,host2")
	env.Set("REDIS_MODE", "cluster")
	env.Set("DATABASE_HOST", "dbhost")
	env.Set("DATABASE_PORT", "5432")
	env.Set("DATABASE_USERNAME", "admin")
	env.Set("DATABASE_PASSWORD", "secret")
	env.Set("DATABASE_NAME", "mydb")

	var cfg Config
	if err := env.Unmarshal(&cfg); err != nil {
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

	env.Unset("DEBUG")
	env.Unset("PORT")
	env.Unset("REDIS_HOST")
	env.Unset("REDIS_MODE")
	env.Unset("DATABASE_HOST")
	env.Unset("DATABASE_PORT")
	env.Unset("DATABASE_USERNAME")
	env.Unset("DATABASE_PASSWORD")
	env.Unset("DATABASE_NAME")
}
