# env

[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/env.svg)](https://pkg.go.dev/github.com/syntaqx/env)
[![codecov](https://codecov.io/gh/syntaqx/env/graph/badge.svg?token=m4bBKy3UG3)](https://codecov.io/gh/syntaqx/env)
[![Go Report Card](https://goreportcard.com/badge/github.com/syntaqx/env)](https://goreportcard.com/report/github.com/syntaqx/env)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

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

## Basic Usage

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

## Unmarshal to Struct

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
	Hosts []string `env:"HOSTS,default=[localhost,localhost2]`
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

### Defaults from Code

You may define default values also in your code by initializing your struct data
before it's populated by `env.Unmarshal`. However, default values defined as
struct tags will take precedence over the ones defined in code.

```go
type Config struct {
    Username string `env:"USERNAME,default=admin"`
    Password string `env:"PASSWORD"`
}

cfg := Config{
    Username: "test",
    Password: "password123",
}

if err := env.Unmarshal(&cfg); err != nil {
    log.Fatalf("Error unmarshalling config: %v", err)
}

// { Username: "admin", Password: "password123" }
```

### From file

The `file` tag option can be used to indicate that the value of the variable
should be loaded from a file. The path of the file given by the value of the
variable.

```bash
echo "password123" > /run/secrets/password
```

```go
type Config struct {
    Username string `env:"USERNAME"`
    Password string `env:"PASSWORD,file"`
}

cfg := Config{
    Username: "test",
    Password: "/run/secrets/password",
}

if err := env.Unmarshal(&cfg); err != nil {
    log.Fatalf("Error unmarshalling config: %v", err)
}

// { "Username": "test", "Password": "password123" }
```

### Expand variables

The `expand` tag option can be used to indicate that the value of the variable
should be expanded (in either `${var}` or `$var` format) before being set.

```go
type Config struct {
    Username string `env:"USERNAME,expand"`
    Password string `env:"PASSWORD,expand"`
}
```

This works great with the `default` tag option:

```go
type Config struct {
    Address string `env:"ADDRESS,expand,default=${HOST}:${PORT}"`
}
```

Which results in:

```bash
HOST=localhost PORT=8080 go run main.go
{Address:localhost:8080}
```

Additionally, default values can be referenced from other struct fields.
Allowing you to chain default values rather than falling back to an empty value
when an environment variable is not set:

```go
type Config struct {
    Host string `env:"HOST,default=localhost"`
    Port string `env:"PORT,default=8080"`
    Address string `env:"ADDRESS,expand,default=${HOST}:${PORT}"`
}
```

Which results in:

```bash
go run main.go
{Host:localhost Port:8080 Address:localhost:8080}
```

## Contributing

Feel free to open issues or contribute to the project. Contributions are always
welcome!

## License

This project is licensed under the MIT license.
