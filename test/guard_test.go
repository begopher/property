package test

import (
	"fmt"
	"github.com/begopher/property"
	"testing"
)

func Test_func_Guard_panic_when_rule_is_nil(t *testing.T) {
	anyProperty := delegate[string]{t: t}
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("passing nil to argument 'constraint' must cause panic")
		}
		expected := "property.Guard: cannot be created from nil constraint"
		if got != expected {
			t.Errorf("Expected message is (%v) got (%v)", expected, got)
		}
	}()
	property.Guard[string](nil, anyProperty)
}

func Test_func_Guard_panic_when_property_is_nil(t *testing.T) {
	anyRule := rule[int]{}
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("passing nil to argument 'property' must cause panic")
		}
		expected := "property.Guard: cannot be created from nil property"
		if got != expected {
			t.Errorf("Expected message is (%v) got (%v)", expected, got)
		}
	}()
	property.Guard[int](anyRule, nil)
}

func Test_guard_Change_delegate_to_underlying_rule_without_mutation(t *testing.T) {
	table := []string{"go", "golang"}
	for _, expected := range table {
		rule := rule[string]{
			evaluate: func(got string) error {
				if got != expected {
					t.Errorf("expected value is (%v) got (%v)", expected, got)
				}
				return nil
			},
		}
		delegation := delegate[string]{
			change: func(string) error { return nil },
		}
		evaluation := property.Guard[string](rule, delegation)
		evaluation.Change(expected)
	}
}

func Test_guard_Change_rule_violation_prevents_delegation_to_underlying_property(t *testing.T) {
	table := []string{"any", "value"}
	for _, data := range table {
		rule := rule[string]{
			evaluate: func(string) error {
				return fmt.Errorf("Any error")
			},
		}
		delegation := delegate[string]{
			change: func(string) error {
				t.Errorf("delegation did occur to underlying property")
				return nil
			},
		}
		eval := property.Guard[string](rule, delegation)
		eval.Change(data)
	}
}

func Test_guard_Change_delegates_to_underlying_property_without_mutation(t *testing.T) {
	table := []string{"go", "golang"}
	for _, expected := range table {
		rule := rule[string]{
			evaluate: func(string) error {
				return nil
			},
		}
		delegation := delegate[string]{
			change: func(got string) error {
				if got != expected {
					t.Errorf("expected value is (%v) got (%v)", expected, got)
				}
				return nil
			},
		}
		eval := property.Guard[string](rule, delegation)
		eval.Change(expected)
	}
}

func Test_guard_Change_returns_error_from_underlying_property_without_mutation(t *testing.T) {
	table := []error{
		nil,
		fmt.Errorf("Any custom error"),
	}
	for _, expected := range table {
		rule := rule[int]{
			evaluate: func(int) error {
				return nil
			},
		}
		delegation := delegate[int]{
			t: t,
			change: func(int) error {
				return expected
			},
		}
		eval := property.Guard[int](rule, delegation)
		var any int
		got := eval.Change(any)
		if got != expected {
			t.Errorf("expected error is (%v) got (%v)", expected, got)
		}
	}
}

func Test_guard_Value_delegates_to_underlying_property(t *testing.T) {
	rule := rule[string]{
		evaluate: func(string) error {
			return nil
		},
	}
	var delegated bool
	delegation := delegate[string]{
		t: t,
		value: func() (string, error) {
			delegated = true
			return "", nil
		},
	}
	eval := property.Guard[string](rule, delegation)
	eval.Value()
	if !delegated {
		t.Error("Delegation did not occur")
	}
}

func Test_guard_Value_returns_value_without_mutation(t *testing.T) {
	table := []struct {
		value string
		err   error
	}{
		{"Go", nil},
		{"any", fmt.Errorf("Any error")},
	}
	for _, data := range table {
		rule := rule[string]{
			evaluate: func(string) error {
				return nil
			},
		}
		delegation := delegate[string]{
			t: t,
			value: func() (string, error) {
				return data.value, data.err
			},
		}
		eval := property.Guard[string](rule, delegation)
		got, _ := eval.Value()
		expected := data.value
		if got != expected {
			t.Errorf("expected value is (%v) got (%v)", expected, got)
		}
	}
}

func Test_guard_Value_returns_error_without_mutation(t *testing.T) {
	table := []struct {
		value string
		err   error
	}{
		{"Go", nil},
		{"any", fmt.Errorf("Any error")},
	}
	for _, data := range table {
		rule := rule[string]{
			evaluate: func(string) error {
				return nil
			},
		}
		delegation := delegate[string]{
			t: t,
			value: func() (string, error) {
				return data.value, data.err
			},
		}
		eval := property.Guard[string](rule, delegation)
		_, got := eval.Value()
		expected := data.err
		if got != expected {
			t.Errorf("expected error is (%v) got (%v)", expected, got)
		}
	}
}
