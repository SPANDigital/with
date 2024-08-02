// SPDX-License-Identifier: MIT
//go:build samples

package ipstack

import (
	"github.com/spandigital/with"
	"reflect"
	"testing"
)

func TestNewIPStack(t *testing.T) {
	type args struct {
		withOptions []with.Func[options]
	}
	tests := []struct {
		name           string
		args           args
		wantNewIpStack *ipstack
		wantErr        bool
	}{
		{
			name: "failure",
			args: args{
				withOptions: []with.Func[options]{},
			},
			wantNewIpStack: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewIpStack, err := NewIPStack(tt.args.withOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIPStack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewIpStack, tt.wantNewIpStack) {
				t.Errorf("NewIPStack() gotNewIpStack = %v, want %v", gotNewIpStack, tt.wantNewIpStack)
			}
		})
	}
}
