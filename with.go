// SPDX-License-Identifier: MIT

// Package with provides a generic implementation of the Functional Options Pattern.
// It allows for flexible configuration of types through functional options while
// supporting defaults, validation, and struct-based initialization.
//
// Example usage:
//
//	type ServerOptions struct {
//	    Host string
//	    Port int
//	}
//
//	func (o *ServerOptions) SetDefaults() {
//	    o.Host = "localhost"
//	    o.Port = 8080
//	}
//
//	func (o *ServerOptions) Validate() error {
//	    if o.Port < 1 || o.Port > 65535 {
//	        return errors.New("invalid port")
//	    }
//	    return nil
//	}
//
//	func WithHost(host string) with.Func[ServerOptions] {
//	    return func(o *ServerOptions) error {
//	        o.Host = host
//	        return nil
//	    }
//	}
//
//	func NewServer(opts ...with.Func[ServerOptions]) (*Server, error) {
//	    o := &ServerOptions{}
//	    if err := with.DefaultThenAddWith(o, opts); err != nil {
//	        return nil, err
//	    }
//	    return &Server{host: o.Host, port: o.Port}, nil
//	}
//
//	// Usage:
//	server, err := NewServer(WithHost("example.com"))
package with

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Defaulted is an optional interface that types can implement to provide default values.
// If the options type implements this interface, SetDefaults will be called by
// DefaultThenAddWith before applying any functional options.
//
// Example:
//
//	type Options struct {
//	    Timeout time.Duration
//	}
//
//	func (o *Options) SetDefaults() {
//	    o.Timeout = 30 * time.Second
//	}
type Defaulted interface {
	SetDefaults()
}

// Validated is an optional interface that types can implement to validate their state.
// If the options type implements this interface, Validate will be called by AddWith
// and DefaultThenAddWith after all functional options have been applied.
//
// Example:
//
//	type Options struct {
//	    Port int
//	}
//
//	func (o *Options) Validate() error {
//	    if o.Port < 1 || o.Port > 65535 {
//	        return errors.New("invalid port range")
//	    }
//	    return nil
//	}
type Validated interface {
	Validate() error
}

// Func is a generic functional option type that modifies options of type O.
// It returns an error if the option cannot be applied.
//
// Example:
//
//	func WithPort(port int) with.Func[ServerOptions] {
//	    return func(o *ServerOptions) error {
//	        if port < 1 || port > 65535 {
//	            return errors.New("invalid port")
//	        }
//	        o.Port = port
//	        return nil
//	    }
//	}
type Func[O any] func(options *O) (err error)

// Nop returns a no-operation functional option that does nothing and always succeeds.
// This can be useful for conditional option application.
//
// Example:
//
//	opts := []with.Func[ServerOptions]{
//	    WithHost("localhost"),
//	}
//	if enableCustomPort {
//	    opts = append(opts, WithPort(9000))
//	} else {
//	    opts = append(opts, with.Nop[ServerOptions]())
//	}
func Nop[O any]() Func[O] {
	return func(options *O) error {
		return nil
	}
}

// DefaultThenAddWith first applies defaults (if the type implements Defaulted),
// then applies the provided functional options, and finally validates (if the type
// implements Validated).
//
// This is the recommended function for constructors that want to start with sensible
// defaults and allow callers to override them.
//
// If any functional option returns an error, processing stops immediately and that
// error is returned with additional context about which option failed.
//
// If validation fails, the validation error is returned.
//
// Example:
//
//	func NewServer(opts ...with.Func[ServerOptions]) (*Server, error) {
//	    o := &ServerOptions{}
//	    if err := with.DefaultThenAddWith(o, opts); err != nil {
//	        return nil, err
//	    }
//	    return &Server{options: o}, nil
//	}
func DefaultThenAddWith[O any](options *O, withOptions []Func[O]) (err error) {
	if i, ok := any(options).(Defaulted); ok {
		i.SetDefaults()
	}
	return AddWith(options, withOptions)
}

// AddWith applies the provided functional options to the options struct, then validates
// (if the type implements Validated).
//
// This is useful when you want to start with a pre-configured options struct or when
// you want to apply additional options to an existing configuration.
//
// If any functional option returns an error, processing stops immediately and that
// error is returned with additional context about which option failed.
//
// If validation fails, the validation error is returned.
//
// Example:
//
//	func NewServerFromOptions(baseOpts *ServerOptions, opts ...with.Func[ServerOptions]) (*Server, error) {
//	    if err := with.AddWith(baseOpts, opts); err != nil {
//	        return nil, err
//	    }
//	    return &Server{options: baseOpts}, nil
//	}
func AddWith[O any](options *O, withOptions []Func[O]) (err error) {
	for _, option := range withOptions {
		if err = option(options); err != nil {
			frame, _ := runtime.CallersFrames([]uintptr{reflect.ValueOf(option).Pointer()}).Next()
			withNames := strings.Split(frame.Function, ".")
			err = fmt.Errorf("cannot apply %v: %w", withNames[len(withNames)-2], err)
			return
		}
	}
	if v, ok := interface{}(options).(Validated); ok {
		if err = v.Validate(); err != nil {
			err = fmt.Errorf("validation failed: %w", err)
		}
	}
	return
}

// MustDefaultThenAddWith is like DefaultThenAddWith but panics if an error occurs.
// This is useful for initialization code where errors should be fatal.
//
// Example:
//
//	var server = NewServerMust(
//	    with.MustDefaultThenAddWith(&ServerOptions{},
//	        []with.Func[ServerOptions]{WithHost("localhost")}))
//
//	func NewServerMust(opts *ServerOptions) *Server {
//	    return &Server{options: opts}
//	}
func MustDefaultThenAddWith[O any](options *O, withOptions []Func[O]) *O {
	if err := DefaultThenAddWith(options, withOptions); err != nil {
		panic(err)
	}
	return options
}

// MustAddWith is like AddWith but panics if an error occurs.
// This is useful for initialization code where errors should be fatal.
//
// Example:
//
//	var defaultOptions = with.MustAddWith(&ServerOptions{Host: "localhost"},
//	    []with.Func[ServerOptions]{WithPort(8080)})
func MustAddWith[O any](options *O, withOptions []Func[O]) *O {
	if err := AddWith(options, withOptions); err != nil {
		panic(err)
	}
	return options
}
