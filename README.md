# with

**A type-safe, generic implementation of the Functional Options Pattern for Go**

[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/SPANDigital/with)
![Develop Go Action Workflow Status](https://img.shields.io/github/actions/workflow/status/spandigital/with/go.yml?branch=develop&label=develop)
![Main Go Action Workflow Status](https://img.shields.io/github/actions/workflow/status/spandigital/with/go.yml?branch=main&label=main)
![Release status](https://img.shields.io/github/v/tag/SPANDigital/with)

## What is this?

The `with` package makes it easy to create clean, flexible APIs in Go using the [Functional Options Pattern](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html). Originally introduced by Rob Pike and popularized by [Dave Cheney](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis), this pattern allows you to write constructors that are:

- **Easy to use** - Simple cases stay simple, complex cases are possible
- **Future-proof** - Add new options without breaking existing code
- **Type-safe** - Leverage Go's generics for compile-time safety
- **Self-documenting** - Options are explicit and readable

## Why use this package?

Instead of writing boilerplate for each type, `with` provides generic helpers that handle:

✅ Setting default values
✅ Applying functional options
✅ Validating the final configuration
✅ Clear error messages when something goes wrong

## Quick Start

### Installation

```bash
go get github.com/spandigital/with
```

### Simple Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/spandigital/with"
)

// 1. Define your options struct
type ServerOptions struct {
    Host    string
    Port    int
    Timeout time.Duration
}

// 2. (Optional) Add defaults
func (o *ServerOptions) SetDefaults() {
    o.Host = "localhost"
    o.Port = 8080
    o.Timeout = 30 * time.Second
}

// 3. (Optional) Add validation
func (o *ServerOptions) Validate() error {
    if o.Port < 1 || o.Port > 65535 {
        return fmt.Errorf("port must be between 1-65535")
    }
    return nil
}

// 4. Create option functions
func WithHost(host string) with.Func[ServerOptions] {
    return func(o *ServerOptions) error {
        o.Host = host
        return nil
    }
}

func WithPort(port int) with.Func[ServerOptions] {
    return func(o *ServerOptions) error {
        if port < 1 || port > 65535 {
            return fmt.Errorf("invalid port: %d", port)
        }
        o.Port = port
        return nil
    }
}

// 5. Use in your constructor
func NewServer(opts ...with.Func[ServerOptions]) (*Server, error) {
    o := &ServerOptions{}
    if err := with.DefaultThenAddWith(o, opts); err != nil {
        return nil, err
    }
    return &Server{options: o}, nil
}

// Usage - clean and readable!
func main() {
    // Use defaults
    server1, _ := NewServer()

    // Override specific options
    server2, _ := NewServer(
        WithHost("0.0.0.0"),
        WithPort(3000),
    )
}
```

## Features

### 🎯 Core Capabilities

- **Generic Type Safety** - Uses Go 1.18+ generics for type-safe option functions
- **Default Values** - Implement `SetDefaults()` to provide sensible defaults
- **Validation** - Implement `Validate()` for comprehensive validation after configuration
- **Error Handling** - Clear error messages with context about which option failed
- **Flexible Usage** - Support both functional options and struct-based initialization
- **Must Variants** - Panic-on-error variants for initialization code (`MustAddWith`, `MustDefaultThenAddWith`)

### 🚀 Two Ways to Configure

**1. Start with defaults, override what you need:**

```go
func NewServer(opts ...with.Func[ServerOptions]) (*Server, error) {
    o := &ServerOptions{}
    if err := with.DefaultThenAddWith(o, opts); err != nil {
        return nil, err
    }
    return &Server{options: o}, nil
}
```

**2. Start with a struct, apply additional options:**

```go
func NewServerFromConfig(config *ServerOptions, opts ...with.Func[ServerOptions]) (*Server, error) {
    if err := with.AddWith(config, opts); err != nil {
        return nil, err
    }
    return &Server{options: config}, nil
}
```

## API Reference

### Main Functions

| Function | Description |
|----------|-------------|
| `with.DefaultThenAddWith(opts, funcs)` | Apply defaults, then options, then validate |
| `with.AddWith(opts, funcs)` | Apply options and validate (no defaults) |
| `with.MustDefaultThenAddWith(opts, funcs)` | Like `DefaultThenAddWith` but panics on error |
| `with.MustAddWith(opts, funcs)` | Like `AddWith` but panics on error |
| `with.Nop[T]()` | Returns a no-op option function |

### Interfaces

| Interface | Method | Description |
|-----------|--------|-------------|
| `Defaulted` | `SetDefaults()` | Called to set default values |
| `Validated` | `Validate() error` | Called to validate the final configuration |

## Complete Example

See the [samples/server](./samples/server) directory for a working example. Here's a taste:

```go
// Create with defaults
server, _ := NewServer()

// Override specific options
server, _ := NewServer(
    WithHost("0.0.0.0"),
    WithPort(3000),
    WithTimeout(60 * time.Second),
)

// Mix struct initialization with options
server, _ := NewServerFromOptions(
    &Options{Host: "localhost", Port: 8080},
    WithTimeout(30 * time.Second),
)

```

## Why Functional Options?

The Functional Options Pattern, introduced by Rob Pike in his 2014 post [Self-referential functions and the design of options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html), solves a common problem in API design: how to make constructors that are both simple to use and flexible to extend.

### The Problem

Traditional approaches have drawbacks:

```go
// 😞 Too many parameters - hard to remember order
NewServer("localhost", 8080, 30*time.Second, true, false, "INFO")

// 😞 Config struct - requires nil checks, verbose for simple cases
NewServer(&Config{Host: "localhost", Port: 8080, ...})

// 😞 Builder pattern - too much ceremony
NewServer().WithHost("localhost").WithPort(8080).Build()
```

### The Solution

Functional options provide a clean middle ground:

```go
// 😊 Clean, readable, self-documenting
server, _ := NewServer(
    WithHost("localhost"),
    WithPort(8080),
)
```

**Benefits:**
- ✅ **Backward compatible** - Adding new options doesn't break existing code
- ✅ **Simple by default** - `NewServer()` works with zero configuration
- ✅ **Self-documenting** - Options are explicit: `WithTimeout(30*time.Second)`
- ✅ **Flexible** - Complex configurations are just as easy as simple ones
- ✅ **Type-safe** - Compiler catches mistakes at build time

### Learn More

- 📖 [Rob Pike's original post](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html) - The genesis of the pattern
- 📖 [Dave Cheney's guide](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) - Popularized the approach
- 📖 [Uber's Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md#functional-options) - Best practices from industry

## Performance

The pattern is lightweight with minimal overhead:

```
BenchmarkAddWith_SingleOption-12        91,654,365     12.82 ns/op    24 B/op    1 allocs/op
BenchmarkAddWith_FourOptions-12         34,602,283     34.55 ns/op    72 B/op    4 allocs/op
BenchmarkDirectInit-12              1,000,000,000      0.22 ns/op     0 B/op    0 allocs/op
```

Each option adds ~10-15ns. For configuration code that runs once at startup, this overhead is negligible.

## Contributing

We welcome contributions! See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](./LICENSE) for details.

---

**Made with ❤️ by [SPAN Digital](https://spandigital.com)**


