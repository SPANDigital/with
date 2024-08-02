//go:build samples

package server

import (
	"github.com/spandigital/with"
	"reflect"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	type args struct {
		withOptions []with.Func[options]
	}
	tests := []struct {
		name          string
		args          args
		wantNewServer *server
		wantErr       bool
	}{
		{
			name: "success",
			args: args{
				withOptions: []with.Func[options]{
					WithHost("localhost"),
					WithPort(1000),
					WithTimeout(3 * time.Second),
				},
			},
			wantNewServer: &server{
				host:    "localhost",
				port:    1000,
				timeout: 3 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				withOptions: []with.Func[options]{},
			},
			wantNewServer: nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewServer, err := NewServer(tt.args.withOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewServer, tt.wantNewServer) {
				t.Errorf("NewServer() gotNewServer = %v, want %v", gotNewServer, tt.wantNewServer)
			}
		})
	}
}
