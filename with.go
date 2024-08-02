package with

import (
	"fmt"
)

type Func[O any] func(options *O) (err error)

func Nop[O any]() Func[O] {
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

func Build[O any](initialOptions *O, validateFunc Func[O], withOptions ...Func[O]) (builtOptions *O, err error) {
	builtOptions = initialOptions
	if err = Compose(builtOptions, withOptions...); err != nil {
		return
	}
	if validateFunc != nil {
		err = validateFunc(builtOptions)
	}
	return
}

func OnCondition[O any](condition bool, withOptions ...Func[O]) Func[O] {
	if condition {
		return func(options *O) error {
			return Compose(options, withOptions...)
		}
	}
	return Nop[O]()
}
