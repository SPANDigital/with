// SPDX-License-Identifier: MIT
//go:build samples

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
		err = errors.New("Host is required")
	case o.Port == 0:
		err = errors.New("Port is required")
	case !(o.Port > 0 && o.Port < 65535):
		err = errors.New("Port must be between 1 and 65535")
	case o.Timeout == 0:
		err = errors.New("Timeout is required")
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
			return errors.New("Port is required")
		case !(port > 0 && port < 65535):
			return errors.New("Port must be between 1 and 65535")
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
