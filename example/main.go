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
	Host []string  `end:"REDIS_HOST|REDIS_HOSTS,default=localhost:6379"`
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
	Database DatabaseConfig `env:"DATABASE"`
	Redis    RedisConfig
}

func main() {
	var cfg Config

	// Set example environment variables
	_ = env.Set("DEBUG", "true")
	_ = env.Set("PORT", "9090")
	_ = env.Set("REDIS_HOST", "host1,host2")
	_ = env.Set("REDIS_MODE", "cluster")
	_ = env.Set("DATABASE_HOST", "dbhost")
	_ = env.Set("DATABASE_PORT", "5432")
	_ = env.Set("DATABASE_USERNAME", "admin")
	_ = env.Set("DATABASE_PASSWORD", "secret")
	_ = env.Set("DATABASE_NAME", "mydb")

	if err := env.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	fmt.Printf("Config: %+v\n", cfg)
}
