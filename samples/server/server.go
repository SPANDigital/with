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

func (s *server) Run() {
	fmt.Printf("server listening on %s:%d\n", s.host, s.port)
}
