// SPDX-License-Identifier: MIT
package with

import (
	"testing"
)

// Benchmark types
type benchOptions struct {
	field1 string
	field2 int
	field3 bool
	field4 float64
}

type benchOptionsWithDefaults struct {
	field1 string
	field2 int
	field3 bool
	field4 float64
}

func (o *benchOptionsWithDefaults) SetDefaults() {
	o.field1 = "default"
	o.field2 = 42
	o.field3 = true
	o.field4 = 3.14
}

type benchOptionsWithValidation struct {
	field1 string
	field2 int
	field3 bool
	field4 float64
}

func (o *benchOptionsWithValidation) Validate() error {
	// Minimal validation to measure overhead
	if o.field2 < 0 {
		return nil
	}
	return nil
}

// Benchmark option functions
func benchWithField1(val string) Func[benchOptions] {
	return func(o *benchOptions) error {
		o.field1 = val
		return nil
	}
}

func benchWithField2(val int) Func[benchOptions] {
	return func(o *benchOptions) error {
		o.field2 = val
		return nil
	}
}

func benchWithField3(val bool) Func[benchOptions] {
	return func(o *benchOptions) error {
		o.field3 = val
		return nil
	}
}

func benchWithField4(val float64) Func[benchOptions] {
	return func(o *benchOptions) error {
		o.field4 = val
		return nil
	}
}

func benchWithField1Defaults(val string) Func[benchOptionsWithDefaults] {
	return func(o *benchOptionsWithDefaults) error {
		o.field1 = val
		return nil
	}
}

func benchWithField2Defaults(val int) Func[benchOptionsWithDefaults] {
	return func(o *benchOptionsWithDefaults) error {
		o.field2 = val
		return nil
	}
}

func benchWithField1Validation(val string) Func[benchOptionsWithValidation] {
	return func(o *benchOptionsWithValidation) error {
		o.field1 = val
		return nil
	}
}

func benchWithField2Validation(val int) Func[benchOptionsWithValidation] {
	return func(o *benchOptionsWithValidation) error {
		o.field2 = val
		return nil
	}
}

// Benchmarks
func BenchmarkAddWith_NoOptions(b *testing.B) {
	opts := &benchOptions{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AddWith(opts, []Func[benchOptions]{})
	}
}

func BenchmarkAddWith_SingleOption(b *testing.B) {
	opts := &benchOptions{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AddWith(opts, []Func[benchOptions]{
			benchWithField1("test"),
		})
	}
}

func BenchmarkAddWith_FourOptions(b *testing.B) {
	opts := &benchOptions{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AddWith(opts, []Func[benchOptions]{
			benchWithField1("test"),
			benchWithField2(123),
			benchWithField3(true),
			benchWithField4(3.14),
		})
	}
}

func BenchmarkDefaultThenAddWith_NoOptions(b *testing.B) {
	opts := &benchOptionsWithDefaults{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DefaultThenAddWith(opts, []Func[benchOptionsWithDefaults]{})
	}
}

func BenchmarkDefaultThenAddWith_TwoOptions(b *testing.B) {
	opts := &benchOptionsWithDefaults{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DefaultThenAddWith(opts, []Func[benchOptionsWithDefaults]{
			benchWithField1Defaults("custom"),
			benchWithField2Defaults(999),
		})
	}
}

func BenchmarkAddWith_WithValidation(b *testing.B) {
	opts := &benchOptionsWithValidation{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AddWith(opts, []Func[benchOptionsWithValidation]{
			benchWithField1Validation("test"),
			benchWithField2Validation(123),
		})
	}
}

func BenchmarkNop(b *testing.B) {
	opts := &benchOptions{}
	nopFunc := Nop[benchOptions]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nopFunc(opts)
	}
}

func BenchmarkMustAddWith_SingleOption(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opts := &benchOptions{}
		_ = MustAddWith(opts, []Func[benchOptions]{
			benchWithField1("test"),
		})
	}
}

func BenchmarkMustDefaultThenAddWith_NoOptions(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opts := &benchOptionsWithDefaults{}
		_ = MustDefaultThenAddWith(opts, []Func[benchOptionsWithDefaults]{})
	}
}

// Benchmark direct struct initialization for comparison
func BenchmarkDirectInit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &benchOptions{
			field1: "test",
			field2: 123,
			field3: true,
			field4: 3.14,
		}
	}
}
