package main

import (
	"fmt"
	"log"

	"github.com/syntaqx/env"
)

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
	Port     string         `env:"PORT,default=8080"`
	Roles    []string       `env:"ROLES,default=[admin,editor]"`
	Database DatabaseConfig `env:"DATABASE"`
	Redis    RedisConfig
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
		"PORT":              "9090",
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
}
