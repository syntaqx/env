# env

[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/env.svg)](https://pkg.go.dev/github.com/syntaqx/env)
[![codecov](https://codecov.io/gh/syntaqx/env/graph/badge.svg?token=m4bBKy3UG3)](https://codecov.io/gh/syntaqx/env)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/env)](https://goreportcard.com/report/github.com/syntaqx/env)

`env` is an environment variable utility package for Go.

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/syntaqx/env"
)

func main() {
    port := env.GetWithFallback("PORT", "8080")
    fmt.Printf("Port: %s\n", port)

    // Assuming the value of HOSTS is a comma-separated list of strings
    // Example: some-host:8000,another-host:8000
    hosts := env.GetSliceWithFallback("HOSTS", []string{"fallback-host-1:8000", "fallback-host-2:8000"})
    fmt.Printf("Hosts: %v\n", hosts)
}
```

### Unmarshal Environment Variables into a Struct

The `Unmarshal` function allows you to load environment variables into a struct
based on struct tags. You can use `default` or `fallback` for fallback values
and `required` to enforce that an environment variable must be set.

```go
package main

import (
	"fmt"
	"log"

	"github.com/syntaqx/env"
)

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST,default=localhost"`
	Port     int    `env:"DATABASE_PORT|DB_PORT,fallback=3306"`
	Username string `env:"DATABASE_USERNAME,default=root"`
	Password string `env:"DATABASE_PASSWORD,required"`
	Database string `env:"DATABASE_NAME"`
}

type Config struct {
	Debug    bool           `env:"DEBUG"`
	Port     string         `env:"PORT,default=8080"`
    Database DatabaseConfig
}

func main() {
	var cfg Config

	// Set example environment variables
	_ = env.Set("DEBUG", "true")
	_ = env.Set("PORT", "9090")
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
```

### Nested Struct Prefixes

You can use nested prefixes to group environment variables. For example, you can

```go
type DatabaseConfig struct {
    Host     string `env:"HOST,default=localhost"`
    Port     int    `env:"PORT,fallback=3306"`
    Username string `env:"USERNAME,default=root"`
    Password string `env:"PASSWORD,required"`
    Database string `env:"NAME"`
}

type Config struct {
    Debug    bool           `env:"DEBUG"`
    Port     string         `env:"PORT,default=8080"`
    Database DatabaseConfig `env:"DATABASE"`
}
```

While will prefix the environment variables in the `DatabaseConfig` struct
with `DATABASE_` as they are looked for. This allows you to reuse the same
struct in multiple places without having to worry about conflicting environment
variables.

```go
type Config struct {
    ReadDatabase DatabaseConfig  `env:"WRITE_DATABASE"`
    WriteDatabase DatabaseConfig `env:"READ_DATABASE"`
}
```
