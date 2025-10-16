// SPDX-License-Identifier: MIT
package with

import (
	"errors"
	"testing"
)

// Test types for testing
type testOptions struct {
	value  string
	number int
}

type testOptionsWithDefaults struct {
	value  string
	number int
}

func (o *testOptionsWithDefaults) SetDefaults() {
	o.value = "default"
	o.number = 42
}

type testOptionsWithValidation struct {
	value  string
	number int
}

func (o *testOptionsWithValidation) Validate() error {
	if o.value == "" {
		return errors.New("value is required")
	}
	if o.number < 0 {
		return errors.New("number must be non-negative")
	}
	return nil
}

type testOptionsWithBoth struct {
	value  string
	number int
}

func (o *testOptionsWithBoth) SetDefaults() {
	o.value = "default"
	o.number = 42
}

func (o *testOptionsWithBoth) Validate() error {
	if o.value == "" {
		return errors.New("value is required")
	}
	if o.number < 0 {
		return errors.New("number must be non-negative")
	}
	return nil
}

// Helper With functions for testing
func withValue(val string) Func[testOptions] {
	return func(o *testOptions) error {
		o.value = val
		return nil
	}
}

func withNumber(num int) Func[testOptions] {
	return func(o *testOptions) error {
		o.number = num
		return nil
	}
}

func withErrorOption() Func[testOptions] {
	return func(o *testOptions) error {
		return errors.New("intentional error")
	}
}

func withValueValidated(val string) Func[testOptionsWithValidation] {
	return func(o *testOptionsWithValidation) error {
		o.value = val
		return nil
	}
}

func withNumberValidated(num int) Func[testOptionsWithValidation] {
	return func(o *testOptionsWithValidation) error {
		o.number = num
		return nil
	}
}

func withValueBoth(val string) Func[testOptionsWithBoth] {
	return func(o *testOptionsWithBoth) error {
		o.value = val
		return nil
	}
}

func withNumberBoth(num int) Func[testOptionsWithBoth] {
	return func(o *testOptionsWithBoth) error {
		o.number = num
		return nil
	}
}

func withValueDefaults(val string) Func[testOptionsWithDefaults] {
	return func(o *testOptionsWithDefaults) error {
		o.value = val
		return nil
	}
}

func withNumberDefaults(num int) Func[testOptionsWithDefaults] {
	return func(o *testOptionsWithDefaults) error {
		o.number = num
		return nil
	}
}

func TestNop(t *testing.T) {
	opts := &testOptions{}
	nopFunc := Nop[testOptions]()

	err := nopFunc(opts)
	if err != nil {
		t.Errorf("Nop() returned error: %v", err)
	}

	// Options should remain unchanged (zero values)
	if opts.value != "" || opts.number != 0 {
		t.Errorf("Nop() modified options: got value=%q, number=%d", opts.value, opts.number)
	}
}

func TestAddWith_NoOptions(t *testing.T) {
	opts := &testOptions{}
	err := AddWith(opts, []Func[testOptions]{})

	if err != nil {
		t.Errorf("AddWith() with no options returned error: %v", err)
	}

	if opts.value != "" || opts.number != 0 {
		t.Errorf("AddWith() modified options unexpectedly: got value=%q, number=%d", opts.value, opts.number)
	}
}

func TestAddWith_SingleOption(t *testing.T) {
	opts := &testOptions{}
	err := AddWith(opts, []Func[testOptions]{withValue("test")})

	if err != nil {
		t.Errorf("AddWith() returned error: %v", err)
	}

	if opts.value != "test" {
		t.Errorf("AddWith() did not apply option: got value=%q, want %q", opts.value, "test")
	}
}

func TestAddWith_MultipleOptions(t *testing.T) {
	opts := &testOptions{}
	err := AddWith(opts, []Func[testOptions]{
		withValue("hello"),
		withNumber(123),
	})

	if err != nil {
		t.Errorf("AddWith() returned error: %v", err)
	}

	if opts.value != "hello" {
		t.Errorf("AddWith() value: got %q, want %q", opts.value, "hello")
	}
	if opts.number != 123 {
		t.Errorf("AddWith() number: got %d, want %d", opts.number, 123)
	}
}

func TestAddWith_OptionError(t *testing.T) {
	opts := &testOptions{}
	err := AddWith(opts, []Func[testOptions]{
		withValue("test"),
		withErrorOption(),
		withNumber(123),
	})

	if err == nil {
		t.Fatal("AddWith() should have returned error")
	}

	// Error message should include function name
	if !contains(err.Error(), "cannot apply") {
		t.Errorf("AddWith() error message should contain 'cannot apply': got %q", err.Error())
	}

	// First option should have been applied, but not the third
	if opts.value != "test" {
		t.Errorf("AddWith() should have applied first option: got value=%q", opts.value)
	}
	if opts.number != 0 {
		t.Errorf("AddWith() should not have applied third option: got number=%d", opts.number)
	}
}

func TestAddWith_WithValidation_Success(t *testing.T) {
	opts := &testOptionsWithValidation{}
	err := AddWith(opts, []Func[testOptionsWithValidation]{
		withValueValidated("test"),
		withNumberValidated(10),
	})

	if err != nil {
		t.Errorf("AddWith() returned error: %v", err)
	}

	if opts.value != "test" || opts.number != 10 {
		t.Errorf("AddWith() options: got value=%q, number=%d", opts.value, opts.number)
	}
}

func TestAddWith_WithValidation_Failure(t *testing.T) {
	opts := &testOptionsWithValidation{}
	err := AddWith(opts, []Func[testOptionsWithValidation]{
		withNumberValidated(10),
		// value is empty, should fail validation
	})

	if err == nil {
		t.Fatal("AddWith() should have returned validation error")
	}

	if !contains(err.Error(), "validation failed") {
		t.Errorf("AddWith() should wrap validation errors, got %q", err.Error())
	}

	if !contains(err.Error(), "required") {
		t.Errorf("AddWith() validation error: got %q", err.Error())
	}
}

func TestDefaultThenAddWith_WithDefaults(t *testing.T) {
	opts := &testOptionsWithDefaults{}
	err := DefaultThenAddWith(opts, []Func[testOptionsWithDefaults]{})

	if err != nil {
		t.Errorf("DefaultThenAddWith() returned error: %v", err)
	}

	if opts.value != "default" {
		t.Errorf("DefaultThenAddWith() value: got %q, want %q", opts.value, "default")
	}
	if opts.number != 42 {
		t.Errorf("DefaultThenAddWith() number: got %d, want %d", opts.number, 42)
	}
}

func TestDefaultThenAddWith_OverrideDefaults(t *testing.T) {
	opts := &testOptionsWithDefaults{}
	err := DefaultThenAddWith(opts, []Func[testOptionsWithDefaults]{
		withValueDefaults("custom"),
		withNumberDefaults(100),
	})

	if err != nil {
		t.Errorf("DefaultThenAddWith() returned error: %v", err)
	}

	if opts.value != "custom" {
		t.Errorf("DefaultThenAddWith() value: got %q, want %q", opts.value, "custom")
	}
	if opts.number != 100 {
		t.Errorf("DefaultThenAddWith() number: got %d, want %d", opts.number, 100)
	}
}

func TestDefaultThenAddWith_PartialOverride(t *testing.T) {
	opts := &testOptionsWithDefaults{}
	err := DefaultThenAddWith(opts, []Func[testOptionsWithDefaults]{
		withValueDefaults("custom"),
	})

	if err != nil {
		t.Errorf("DefaultThenAddWith() returned error: %v", err)
	}

	// Value should be overridden, number should be default
	if opts.value != "custom" {
		t.Errorf("DefaultThenAddWith() value: got %q, want %q", opts.value, "custom")
	}
	if opts.number != 42 {
		t.Errorf("DefaultThenAddWith() number: got %d, want %d (default)", opts.number, 42)
	}
}

func TestDefaultThenAddWith_NoDefaultsInterface(t *testing.T) {
	opts := &testOptions{}
	err := DefaultThenAddWith(opts, []Func[testOptions]{
		withValue("test"),
	})

	if err != nil {
		t.Errorf("DefaultThenAddWith() returned error: %v", err)
	}

	// Should work like AddWith when no Defaulted interface
	if opts.value != "test" {
		t.Errorf("DefaultThenAddWith() value: got %q, want %q", opts.value, "test")
	}
	if opts.number != 0 {
		t.Errorf("DefaultThenAddWith() number: got %d, want %d", opts.number, 0)
	}
}

func TestDefaultThenAddWith_WithBoth(t *testing.T) {
	opts := &testOptionsWithBoth{}
	err := DefaultThenAddWith(opts, []Func[testOptionsWithBoth]{
		withValueBoth("custom"),
		withNumberBoth(100),
	})

	if err != nil {
		t.Errorf("DefaultThenAddWith() returned error: %v", err)
	}

	if opts.value != "custom" {
		t.Errorf("DefaultThenAddWith() value: got %q, want %q", opts.value, "custom")
	}
	if opts.number != 100 {
		t.Errorf("DefaultThenAddWith() number: got %d, want %d", opts.number, 100)
	}
}

func TestDefaultThenAddWith_WithBoth_ValidationFailure(t *testing.T) {
	opts := &testOptionsWithBoth{}
	err := DefaultThenAddWith(opts, []Func[testOptionsWithBoth]{
		withNumberBoth(-10), // negative number should fail validation
	})

	if err == nil {
		t.Fatal("DefaultThenAddWith() should have returned validation error")
	}

	if !contains(err.Error(), "validation failed") {
		t.Errorf("DefaultThenAddWith() should wrap validation errors, got %q", err.Error())
	}

	if !contains(err.Error(), "non-negative") {
		t.Errorf("DefaultThenAddWith() validation error: got %q", err.Error())
	}
}

func TestMustDefaultThenAddWith_Success(t *testing.T) {
	opts := &testOptionsWithDefaults{}
	result := MustDefaultThenAddWith(opts, []Func[testOptionsWithDefaults]{
		withValueDefaults("custom"),
	})

	if result != opts {
		t.Error("MustDefaultThenAddWith() should return the same pointer")
	}

	if opts.value != "custom" {
		t.Errorf("MustDefaultThenAddWith() value: got %q, want %q", opts.value, "custom")
	}
	if opts.number != 42 {
		t.Errorf("MustDefaultThenAddWith() number: got %d, want %d", opts.number, 42)
	}
}

func TestMustDefaultThenAddWith_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDefaultThenAddWith() should have panicked")
		}
	}()

	opts := &testOptionsWithBoth{}
	MustDefaultThenAddWith(opts, []Func[testOptionsWithBoth]{
		withNumberBoth(-10), // will fail validation
	})
}

func TestMustAddWith_Success(t *testing.T) {
	opts := &testOptions{}
	result := MustAddWith(opts, []Func[testOptions]{
		withValue("test"),
		withNumber(123),
	})

	if result != opts {
		t.Error("MustAddWith() should return the same pointer")
	}

	if opts.value != "test" {
		t.Errorf("MustAddWith() value: got %q, want %q", opts.value, "test")
	}
	if opts.number != 123 {
		t.Errorf("MustAddWith() number: got %d, want %d", opts.number, 123)
	}
}

func TestMustAddWith_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustAddWith() should have panicked")
		}
	}()

	opts := &testOptions{}
	MustAddWith(opts, []Func[testOptions]{
		withErrorOption(),
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
