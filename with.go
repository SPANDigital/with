// SPDX-License-Identifier: MIT
package with

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type Defaulted interface {
	SetDefaults()
}

type Validated interface {
	Validate() error
}

type Func[O any] func(options *O) (err error)

func Nop[O any]() Func[O] {
	return func(options *O) error {
		return nil
	}
}

func DefaultThenAddWith[O any](options *O, withOptions []Func[O]) (err error) {
	if i, ok := any(options).(Defaulted); ok {
		i.SetDefaults()
	}
	return AddWith(options, withOptions)
}

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
		err = v.Validate()
	}
	return
}
