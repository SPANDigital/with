# with
SPAN Digital's take on the Functional Options Pattern using Generics

[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/SPANDigital/with)
![Develop Go Action Workflow Status](https://img.shields.io/github/actions/workflow/status/spandigital/with/go.yml?branch=develop&label=develop)
![Main Go Action Workflow Status](https://img.shields.io/github/actions/workflow/status/spandigital/with/go.yml?branch=main&label=main)
![Release status](https://img.shields.io/github/v/tag/SPANDigital/with)

### Features
- Supports default options (use ```SetDefaults()``` pointer receiver )
- - Supports validation of options as a whole (use ```Validate()``` pointer receiver)
- Supports validation of individual `With...` functions, by returning an error
- Allows for option to be passed as structs, _you can ignore functional options if you so wish_
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

3. Write constructors using

   To start with defaults, and apply 0..n use ``With...`` functions

   To start with a options str

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

type Options struct {
   Host    string
   Port    int
   Timeout time.Duration
}

func (o *Options) Validate() (err error) {
   switch {
   case o.Host == "":
      err = errors.New("host is required")
   case o.Port == 0:
      err = errors.New("port is required")
   case !(o.Port > 0 && o.Port < 65535):
      err = errors.New("port must be between 1 and 65535")
   case o.Timeout == 0:
      err = errors.New("timeout is required")
   }
   return
}

func WithHost(host string) with.Func[Options] {
   return func(options *Options) (err error) {
      options.Host = host
      return
   }
}

func WithPort(port int) with.Func[Options] {
   return func(options *Options) (err error) {
      switch {
      case port == 0:
         return errors.New("port is required")
      case !(port > 0 && port < 65535):
         return errors.New("port must be between 1 and 65535")
      }
      options.Port = port
      return
   }
}

func WithTimeout(timeout time.Duration) with.Func[Options] {
   return func(options *Options) (err error) {
      options.Timeout = timeout
      return
   }
}

type server struct {
   host    string
   port    int
   timeout time.Duration
}

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

func (s *server) Run() {
   fmt.Printf("server listening on %s:%d\n", s.host, s.port)
}
```


