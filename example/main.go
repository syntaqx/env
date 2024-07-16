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

func main() {
	var cfg Config

	// Set example environment variables
	env.Set("DEBUG", "true")
	env.Set("PORT", "9090")
	env.Set("REDIS_HOST", "host1,host2")
	env.Set("REDIS_MODE", "cluster")
	env.Set("DATABASE_HOST", "dbhost")
	env.Set("DATABASE_PORT", "5432")
	env.Set("DATABASE_USERNAME", "admin")
	env.Set("DATABASE_PASSWORD", "secret")
	env.Set("DATABASE_NAME", "mydb")

	if err := env.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	fmt.Printf("Config: %+v\n", cfg)
}
