package with

import (
	"reflect"
	"testing"
)

type Options struct {
	alpha   string
	beta    string
	gamma   string
	delta   int
	epsilon int
	zeta    bool
}

func WithAlpha(alpha string) Func[Options] {
	return func(options *Options) error {
		options.alpha = alpha
		return nil
	}
}

func WithBeta(beta string) Func[Options] {
	return func(options *Options) error {
		options.beta = beta
		return nil
	}
}

func WithGamma(gamma string) Func[Options] {
	return func(options *Options) error {
		options.gamma = gamma
		return nil
	}
}

func WithDelta(delta int) Func[Options] {
	return func(options *Options) error {
		options.delta = delta
		return nil
	}
}

func WithEpsilon(epsilon int) Func[Options] {
	return func(options *Options) error {
		options.epsilon = epsilon
		return nil
	}
}

func WithZeta(zeta bool) Func[Options] {
	return func(options *Options) error {
		options.zeta = zeta
		return nil
	}
}

func TestBuild(t *testing.T) {
	type args[O any] struct {
		initialOptions *O
		postFunc       Func[O]
		withOptions    []Func[O]
	}
	type testCase[O any] struct {
		name             string
		args             args[O]
		wantBuiltOptions *O
		wantErr          bool
	}
	tests := []testCase[Options]{
		{
			name: "All",
			args: args[Options]{
				initialOptions: &Options{
					zeta: false,
				},
				postFunc: WithZeta(true),
				withOptions: []Func[Options]{
					WithAlpha("alpha"),
					WithBeta("beta"),
					WithGamma("gamma"),
					WithDelta(-1),
					WithEpsilon(-2),
				},
			},
			wantBuiltOptions: &Options{
				alpha:   "alpha",
				beta:    "beta",
				gamma:   "gamma",
				delta:   -1,
				epsilon: -2,
				zeta:    true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuiltOptions, err := Build(tt.args.initialOptions, tt.args.postFunc, tt.args.withOptions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBuiltOptions, tt.wantBuiltOptions) {
				t.Errorf("Build() gotBuiltOptions = %v, want %v", gotBuiltOptions, tt.wantBuiltOptions)
			}
		})
	}
}

func TestCompose(t *testing.T) {
	type args[O Options] struct {
		options     *O
		withOptions []Func[Options]
	}
	type testCase[O any] struct {
		name    string
		args    args[Options]
		wantErr bool
	}
	tests := []testCase[Options]{
		{
			name: "All",
			args: args[Options]{
				options: &Options{},
				withOptions: []Func[Options]{
					WithAlpha("alpha"),
					WithBeta("beta"),
					WithGamma("gamma"),
					WithDelta(-1),
					WithEpsilon(-2),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Compose(tt.args.options, tt.args.withOptions...); (err != nil) != tt.wantErr {
				t.Errorf("Compose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
