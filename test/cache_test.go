package test

import (
	"fmt"
	"testing"

	"github.com/begopher/property"
)

func Test_func_Cache_panic_when_property_is_nil(t *testing.T) {
	var any int
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("passing nil to argument 'property' must cause panic")
		}
		expected := "property.Cache: cannot be created from nil property"
		if got != expected {
			t.Errorf("Expected message is (%v) got (%v)", expected, got)
		}
	}()
	property.Cache[int](any, nil)
}

func Test_cache_Value_returns_what_is_in_cache(t *testing.T) {
	table := []string{"any", "go", "golang"}
	delegate := delegate[string]{t: t}
	for _, expected := range table {
		cache := property.Cache[string](expected, delegate)
		got, err := cache.Value()
		if got != expected {
			t.Errorf("Expected value is (%v) got (%v)", expected, got)
		}
		if err != nil {
			t.Errorf("Expected error is (nil) got (%v)", err)
		}
	}
}

func Test_cache_Value_does_not_delegate_to_underlying_property(t *testing.T) {
	delegate := delegate[int]{
		t: t,
		value: func() (int, error) {
			t.Errorf("delegation did occur")
			var anyValue int
			var anyError error
			return anyValue, anyError
		},
	}
	var any int
	cache := property.Cache[int](any, delegate)
	cache.Value()
}

func Test_cache_Change_delegates_to_underlying_property_without_mutation(t *testing.T) {
	table := []string{"any", "go", "golang"}
	for _, expected := range table {
		var delegated bool
		delegate := delegate[string]{
			t: t,
			change: func(got string) error {
				delegated = true
				if expected != got {
					t.Errorf("Expected value is (%v) got (%v)", expected, got)
				}
				return nil
			},
		}
		var any string
		cache := property.Cache[string](any, delegate)
		cache.Change(expected)
		if !delegated {
			t.Error("delegation to underlying property never occur")
		}
	}
}

func Test_cache_Change_returns_error_of_underlying_property_without_mutation(t *testing.T) {
	table := []error{
		nil,
		fmt.Errorf("custom error 1"),
		fmt.Errorf("custom error 2"),
	}
	for _, err := range table {
		delegate := delegate[int]{
			t: t,
			change: func(int) error {
				return err
			},
		}
		var any int
		cache := property.Cache[int](any, delegate)
		got := cache.Change(any)
		expected := err
		if got != expected {
			t.Errorf("Expected error is (%v) got (%v)", expected, got)
		}
	}
}

func Test_cache_Change_when_underlying_property_does_not_return_error_cache_get_updated(t *testing.T) {
	table := []struct {
		oldValue int
		newValue int
	}{
		{0, 1},
		{2, 3},
		{4, 5},
	}

	delegate := delegate[int]{
		t:      t,
		change: func(int) error { return nil },
	}
	for _, data := range table {
		cache := property.Cache[int](data.oldValue, delegate)
		cache.Change(data.newValue)
		got, _ := cache.Value()
		expected := data.newValue
		if got != expected {
			t.Errorf("Expected value is (%v) got (%v)", expected, got)
		}
	}
}

func Test_cache_Change_when_underlying_property_returns_error_cache_does_not_get_updated(t *testing.T) {
	table := []struct {
		oldValue string
		newValue string
	}{
		{"go", "golang"},
		{"C", "C programming langauge"},
	}
	delegate := delegate[string]{
		t: t,
		change: func(string) error {
			return fmt.Errorf("any custom error")
		},
	}
	for _, data := range table {
		cache := property.Cache[string](data.oldValue, delegate)
		cache.Change(data.newValue)
		got, _ := cache.Value()
		expected := data.oldValue
		if got != expected {
			t.Errorf("Expected value is (%v) got (%v)", expected, got)
		}
	}
}
