package with

import (
	"fmt"
)

type Func[O any] func(options *O) error

func Noop[O any]() Func[O] {
	return func(options *O) error {
		return nil
	}
}

func Compose[O any](options *O, withOptions ...Func[O]) (err error) {
	for _, option := range withOptions {
		if err = option(options); err != nil {
			err = fmt.Errorf("cannot apply %v option to %v: %w", option, options, err)
			return
		}
	}
	return
}

func Build[O any](initialOptions *O, postFunc Func[O], withOptions ...Func[O]) (builtOptions *O, err error) {
	builtOptions = initialOptions
	if err = Compose(builtOptions, withOptions...); err != nil {
		return
	}
	if postFunc != nil {
		err = postFunc(builtOptions)
	}
	return
}

func OnCondition[O any](condition bool, withOptions ...Func[O]) Func[O] {
	if condition {
		return func(options *O) error {
			return Compose(options, withOptions...)
		}
	}
	return Noop[O]()
}
