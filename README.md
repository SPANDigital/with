# with
SPAN Digital's take on the Functional Options Pattern using Generics

![build](https://github.com/spandigital/with/actions/workflows/go.yml/badge.svg)

### Features
- Supports default options
- Supports validation of individual `With...` functions 
- Supports validation of options as a whole 
- Allows for config structs 
- Allows callers to create their own `With...` functions

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
   NB: If you want callers to modify options export options and parameters by capitalizing. 
```go
type Options struct {
   Host    string
   Port    int
   Timeout time.Duration
}
```

2. Optionally, write a SetDefaults function (with.Defaultable interface)

```go
func (o *Options) SetDefaults() {
   o.Timeout = 5 * time.Minutes
}
```

3. Optionally, write a Validate function (with.Validated interface) 

```go
func (o *Options) Validate() (err error) {
   switch {
      case o.Host == "":
        err = errors.New("host is required")
      case o.Port == 0:
        err = errors.New("port is required")
      case !(o.Port > 0 && o.port < 65535):
        err = errors.New("port must be between 1 and 65535")
      case o.Timeout == 0:
        err = errors.New("timeout is required")
      }
   return
}
```

3. Optionally, Write High Order Functions prefixed with ``With...`` which manipulate the options.
   You have an option of returning an error the option is not valid.

```go
func WithHost(host string) with.Func[Options] {
   return func(options *Options) (err error) {
      options.Host = host
      return
   }
}

func WithPort(port int) with.Func[Options] {
   return func(options *Options) (err error) {
      switch {
         case Port == 0:
            return errors.New("port is required")
         case !(Port > 0 && Port < 65535):
            return errors.New("port must be between 1 and 65535")
      }
      options.Port = port
      return
   }
}

func WithTimeout(timeout time.Duration) with.Func[Options] {
   return func(options *options) (err error) {
      options.Timeout = timeout
      return
   }
}
```

3. Using `with.Build` to retrieved your configured options, typically this is used in constructors

   Defaults and validation are configured at this step. If you do not require validation use nil
   for the validation function.

```go
func NewServer(withOptions ...with.Func[Options]) (server *server, err error) {
   o := &Options{}
   if err = with.DefaultThenAddWith(o, withOptions); err == nil {
	   server = newServer(o)
   }
   return
}

func NewServerFromOptions(options *Options, withOptions ...with.Func[Options]) (server *server, err error) {
   if err = with.AddWith(options, withOptions); err == nil {
	   server = newServer(options)
   }
   return
}

func newServer(options *Options) *server {
   return &server{
     host:    options.Host,
     port:    options.Port,
     timeout: options.Timeout,
   }
}
```

4. Use:

```go
server, err := NewServer(
	WithHost("localhost"),
	WithPort(10000),
	WithTimeout(3 * time.Second),
)
```

or

```go
server, err = NewServerFromOptions(
	&Options{
	    Host: "localhost",	
    },   
)
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


