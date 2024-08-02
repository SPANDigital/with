# with
SPANDigital's take on the Functional Options Pattern using Generics

![build](https://github.com/spandigital/with/actions/workflows/go.yml/badge.svg)

### Features
- Use any struct for options via generics
- Supports default options
- Supports validation on individual options as they're set
- Supports validation of options as a unit after they're set
- Supports optional grouping of functions using 

### Installation

```bash
go get github.com/spandigital/with
```

### How to use

1. Import with package

```go
import github.com/spandigital/with
```

1. Define a type to represent your options, typically a struct, but not necessarily.

```go
type options struct {
	host string
	port   int
}
```

2. Write High Order Functions prefixed with ``With...`` which manipulate the options.
   You have an option of returning an error the option is not valid.

```go
func WithHost(host string) with.Func[options] {
    return func(options *options) (err error) {
        options.host = host
        return
    }
}
```

3. Using `with.Build` to retrieved your configured options, typically this is used in constructors

   Defaults and validation are configured at this step. If you do not require validation use nil
   for the validation function.

```go
func NewServer(withOptions ...with.Func[options]) (newServer *server, err error) {
	var builtOptions *options
	if builtOptions, err = with.Build(&options{ // defaults
		host:    "",
		port:    0,
		timeout: 0,
	}, func(options *options) (err error) { // validation
		switch {
		case options.host == "":
			err = errors.New("host is required")
		case options.port == 0:
			err = errors.New("port is required")
		case options.timeout == 0:
			err = errors.New("timeout is required")
		}
		return
	}, withOptions...); err == nil { //options is ready to be used
		newServer = &server{
			host:    builtOptions.host,
			port:    builtOptions.port,
			timeout: builtOptions.timeout,
		}
	}
	return
}
```


### Usage samples

See [/samples](/samples) for usage samples.

```go
package server

import (
	"errors"
	"fmt"
	"github.com/spandigital/with"
	"time"
)

type Server interface {
	Run()
}

type options struct {
	host    string
	port    int
	timeout time.Duration
}

func WithHost(host string) with.Func[options] {
	return func(options *options) (err error) {
		options.host = host
		return
	}
}

func WithPort(port int) with.Func[options] {
	return func(options *options) (err error) {
		options.port = port
		return
	}
}

func WithTimeout(timeout time.Duration) with.Func[options] {
	return func(options *options) (err error) {
		options.timeout = timeout
		return
	}
}

type server struct {
	host    string
	port    int
	timeout time.Duration
}

func NewServer(withOptions ...with.Func[options]) (newServer *server, err error) {
	var newOptions *options
	if newOptions, err = with.Build(&options{ // defaults
		host:    "",
		port:    0,
		timeout: 0,
	}, func(options *options) (err error) { // validation
		switch {
		case options.host == "":
			err = errors.New("host is required")
		case options.port == 0:
			err = errors.New("port is required")
		case options.timeout == 0:
			err = errors.New("timeout is required")
		}
		return
	}, withOptions...); err == nil { //options is ready to be used
		newServer = &server{
			host:    newOptions.host,
			port:    newOptions.port,
			timeout: newOptions.timeout,
		}
	}
	return
}

func (s *server) Run() {
	fmt.Printf("server listening on %s:%d\n", s.host, s.port)
}
```


