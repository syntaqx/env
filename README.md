# env

[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/env.svg)](https://pkg.go.dev/github.com/syntaqx/env)
[![codecov](https://codecov.io/gh/syntaqx/env/graph/badge.svg?token=m4bBKy3UG3)](https://codecov.io/gh/syntaqx/env)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/env)](https://goreportcard.com/report/github.com/syntaqx/env)

`env` is an environment variable utility package for Go. It provides simple
functions to get and set environment variables, including support for
unmarshalling environment variables into structs with support for nested
structures, default values, and required fields.

## Features

- __Basic Get/Set__: Simple functions to get, set, and unset environment variables.
- __Type Conversion__: Functions to get environment variables as different types (int, bool, float).
- __Fallback Values__: Support for fallback values if an environment variable is not set.
- __Unmarshal__: Load environment variables into structs using struct tags.
- __Nested Structs__: Support for nested struct prefixes to group environment variables.

## Installation

```sh
go get github.com/syntaqx/env
```

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
    hosts, err := env.GetStringSliceWithFallback("HOSTS", []string{"fallback-host-1:8000", "fallback-host-2:8000"})
    if err != nil {
        fmt.Printf("Error getting hosts: %v\n", err)
    } else {
        fmt.Printf("Hosts: %v\n", hosts)
    }
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

You can use nested prefixes to group environment variables. This allows you to
reuse the same struct in multiple places without having to worry about
conflicting environment variables.


```go
package main

import (
    "fmt"
    "log"

    "github.com/syntaqx/env"
)

type DatabaseConfig struct {
    Host     string `env:"HOST,default=localhost"`
    Port     int    `env:"PORT,fallback=3306"`
    Username string `env:"USERNAME,default=root"`
    Password string `env:"PASSWORD,required"`
    Database string `env:"NAME"`
}

type Config struct {
    Debug         bool           `env:"DEBUG"`
    Port          string         `env:"PORT,default=8080"`
    ReadDatabase  DatabaseConfig `env:"READ_DATABASE"`
    WriteDatabase DatabaseConfig `env:"WRITE_DATABASE"`
}

func main() {
    var cfg Config

    // Set example environment variables
    _ = env.Set("DEBUG", "true")
    _ = env.Set("PORT", "9090")
    _ = env.Set("READ_DATABASE_HOST", "read-dbhost")
    _ = env.Set("READ_DATABASE_PORT", "5432")
    _ = env.Set("READ_DATABASE_USERNAME", "read-admin")
    _ = env.Set("READ_DATABASE_PASSWORD", "read-secret")
    _ = env.Set("READ_DATABASE_NAME", "read-mydb")
    _ = env.Set("WRITE_DATABASE_HOST", "write-dbhost")
    _ = env.Set("WRITE_DATABASE_PORT", "5432")
    _ = env.Set("WRITE_DATABASE_USERNAME", "write-admin")
    _ = env.Set("WRITE_DATABASE_PASSWORD", "write-secret")
    _ = env.Set("WRITE_DATABASE_NAME", "write-mydb")

    if err := env.Unmarshal(&cfg); err != nil {
        log.Fatalf("Error unmarshalling config: %v", err)
    }

    fmt.Printf("Config: %+v\n", cfg)
}
```

### Slice Types Defaults

When using slice types, if you are declaring a single value as the default you
can use the `default` tag as normal:

```go
type Config struct {
	Hosts []string `env:"HOSTS,default=localhost"`
}
```

However if you want to declare multiple values as the default, you must enclose
the values in square brackets:

```go
type Config struct {
	Hosts []string `env:"HOSTS,default=[localhost,localhost2]
}
```

This is necessary as the pacakge uses commas as a delimiter to split the struct
tag options, and without the square brackets it would split the values into
multiple tags.

```go
type Config struct {
	Hosts []string `env:"HOSTS,default=[localhost,localhost2],required"
}
```

## Contributing

Feel free to open issues or contribute to the project. Contributions are always
welcome!

## License

This project is licensed under the MIT license.
