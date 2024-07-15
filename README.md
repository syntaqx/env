# env

[![Go Reference](https://pkg.go.dev/badge/github.com/syntaqx/env.svg)](https://pkg.go.dev/github.com/syntaqx/env)
[![codecov](https://codecov.io/gh/syntaqx/env/graph/badge.svg?token=m4bBKy3UG3)](https://codecov.io/gh/syntaqx/env)

`env` is an environment variable utility package.

> [!NOTE]
> This project is a work in progress as I build out functionality I require in
> other projects I intend to use this package to support them.

### Usage

```go
package main

import (
    "fmt"
    "github.com/syntaqx/env"
)

func main() {
    if err := env.Load(); err != nil {
        fmt.Printf("failed to load environment variables %v\n", err)
    }

    port := env.GetWithFallback("PORT", "8080")
    fmt.Printf("Port: %s\n", port)

    // Assuming the value of HOSTS is a comma separated list of strings
    // Example: some-host:8000,another-host:8000
    hosts := env.GetSliceWithFallback("HOSTS", []string{"fallback-host-1:8000", "fallback-host-2:8000"})
    fmt.Printf("Hosts: %v\n", hosts)
}
```

### Roadmap

- [ ] Load environment variables into structs with tags
- [ ] Type casting of environment variable values
