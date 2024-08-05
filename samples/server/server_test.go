package server

import (
	"github.com/spandigital/with"
	"reflect"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	type args struct {
		withOptions []with.Func[Options]
	}
	tests := []struct {
		name       string
		args       args
		wantServer *server
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				withOptions: []with.Func[Options]{
					WithHost("localhost"),
					WithPort(8080),
					WithTimeout(3 * time.Second),
				},
			},
			wantServer: &server{
				host:    "localhost",
				port:    8080,
				timeout: 3 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "incomplete details 1",
			args: args{
				withOptions: []with.Func[Options]{
					WithHost("localhost"),
				},
			},
			wantServer: nil,
			wantErr:    true,
		},
		{
			name: "incomplete details 2",
			args: args{
				withOptions: []with.Func[Options]{
					WithHost("localhost"),
					WithPort(8080),
				},
			},
			wantServer: nil,
			wantErr:    true,
		},
		{
			name: "invalid Host",
			args: args{
				withOptions: []with.Func[Options]{
					WithHost(""),
				},
			},
			wantServer: nil,
			wantErr:    true,
		},
		{
			name: "invalid Port",
			args: args{
				withOptions: []with.Func[Options]{
					WithPort(-1),
				},
			},
			wantServer: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServer, err := NewServer(tt.args.withOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotServer, tt.wantServer) {
				t.Errorf("NewServer() gotServer = %v, want %v", gotServer, tt.wantServer)
			}
		})
	}
}

func TestNewServerFromOptions(t *testing.T) {
	type args struct {
		options     *Options
		withOptions []with.Func[Options]
	}
	tests := []struct {
		name       string
		args       args
		wantServer *server
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				options: &Options{
					Host:    "localhost",
					Port:    8080,
					Timeout: 3 * time.Second,
				},
			},
			wantServer: &server{
				host:    "localhost",
				port:    8080,
				timeout: 3 * time.Second,
			},
		},
		{
			name: "incomplete details 1",
			args: args{
				options: &Options{
					Host: "localhost",
					Port: 8080,
				},
			},
			wantServer: nil,
			wantErr:    true,
		},
		{
			name: "incomplete details 2",
			args: args{
				options: &Options{
					Host: "localhost",
				},
			},
			wantServer: nil,
			wantErr:    true,
		},
		{
			name: "cobined success",
			args: args{
				options: &Options{
					Host: "localhost",
					Port: 8080,
				},
				withOptions: []with.Func[Options]{
					WithTimeout(3 * time.Second),
				},
			},
			wantServer: &server{
				host:    "localhost",
				port:    8080,
				timeout: 3 * time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServer, err := NewServerFromOptions(tt.args.options, tt.args.withOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerFromOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotServer, tt.wantServer) {
				t.Errorf("NewServerFromOptions() gotServer = %v, want %v", gotServer, tt.wantServer)
			}
		})
	}
}
