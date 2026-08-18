package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/syntaqx/env"
)

// Port accepts either "8080" or ":8080" and normalizes it by trimming a leading
// ":". It stays a string so it can be passed straight to net.JoinHostPort.
// Because it implements encoding.TextUnmarshaler, env.Unmarshal uses it
// automatically wherever a Port field appears.
type Port string

func (p *Port) UnmarshalText(text []byte) error {
	s := strings.TrimPrefix(strings.TrimSpace(string(text)), ":")
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("invalid port %q: %w", text, err)
	}
	*p = Port(s)
	return nil
}

type RedisMode string

const (
	RedisModeStandalone RedisMode = "standalone"
	RedisModeCluster    RedisMode = "cluster"
)

type RedisConfig struct {
	Host []string  `env:"REDIS_HOST|REDIS_HOSTS,default=localhost:6379"`
	Mode RedisMode `env:"REDIS_MODE,default=standalone"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST,default=localhost"`
	Port     int    `env:"PORT|DB_PORT,fallback=3306"`
	Username string `env:"USERNAME,default=root"`
	Password string `env:"PASSWORD,required"`
	Database string `env:"NAME"`
}

type Config struct {
	Debug    bool           `env:"DEBUG"`
	Host     string         `env:"HOST,default=localhost"`
	Port     Port           `env:"PORT,default=8080"`
	Roles    []string       `env:"ROLES,default=[admin,editor]"`
	Database DatabaseConfig `env:"DATABASE"`
	Redis    RedisConfig
}

// Addr derives the listen address from Host and Port together, rather than
// reading it from a separate environment variable.
func (c Config) Addr() string {
	return net.JoinHostPort(c.Host, string(c.Port))
}

func setEnvVars(vars map[string]string) {
	for key, value := range vars {
		if err := env.Set(key, value); err != nil {
			log.Fatalf("Error setting environment variable %s: %v", key, err)
		}
	}
}

func main() {
	envVars := map[string]string{
		"DEBUG":             "true",
		"PORT":              ":9090",
		"REDIS_HOST":        "host1,host2",
		"REDIS_MODE":        "cluster",
		"DATABASE_HOST":     "dbhost",
		"DATABASE_PORT":     "5432",
		"DATABASE_USERNAME": "admin",
		"DATABASE_PASSWORD": "secret",
		"DATABASE_NAME":     "mydb",
	}

	setEnvVars(envVars)

	var cfg Config
	if err := env.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	fmt.Printf("Config: %+v\n", cfg)
	fmt.Printf("Addr:   %s\n", cfg.Addr())
}
