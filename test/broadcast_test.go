package test

import (
	"fmt"
	"testing"

	"github.com/begopher/property"
)

func Test_func_Broadcast_panic_when_receivers_is_nil(t *testing.T) {
	any := delegate[string]{}
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("passing nil receivers, panic must occur")
		}
		expected := "property.broadcast: cannot be created with zero receivers"
		if got != expected {
			t.Errorf("expected message is (%v) got (%v)", expected, got)
		}
	}()
	property.Broadcast[string](nil, any)
}

func Test_func_Broadcast_panic_when_property_is_nil(t *testing.T) {
	anyReceivers := []func(string){
		func(string) {},
	}
	defer func() {
		got := recover()
		if got == nil {
			t.Fatal("passing nil property, panic must occur")
		}
		expected := "property.broadcast: cannot be created from nil property"
		if got != expected {
			t.Errorf("expected message is (%v) got (%v)", expected, got)
		}
	}()
	property.Broadcast[string](anyReceivers, nil)
}

func Test_broadcast_Change_delegate_to_underlying_property(t *testing.T) {
	var invoked bool
	delegation := delegate[int]{
		t: t,
		change: func(int) error {
			invoked = true
			return nil
		},
	}
	receivers := []func(int){
		func(int) { /* any implementaion */ },
	}
	broadcast := property.Broadcast[int](receivers, delegation)
	var any int
	broadcast.Change(any)
	if !invoked {
		t.Errorf("delegation did not occur to underlying property")
	}
}

func Test_broadcast_Change_delegates_to_underlying_property_without_mutation(t *testing.T) {
	table := []string{"Go", "Golang"}
	for _, expected := range table {
		delegation := delegate[string]{
			t: t,
			change: func(got string) error {
				if got != expected {
					t.Errorf("expected value is (%v) got (%v)", expected, got)
				}
				return nil
			},
		}
		receivers := []func(string){
			func(got string) { /* any implementation */ },
		}
		broadcast := property.Broadcast[string](receivers, delegation)
		broadcast.Change(expected)
	}
}

func Test_broadcast_Change_notifies_all_receivers(t *testing.T) {
	var receiver1, receiver2, receiver3 bool
	delegation := delegate[int]{
		t:      t,
		change: func(int) error { return nil },
	}
	receivers := []func(int){
		func(int) { receiver1 = true },
		func(int) { receiver2 = true },
		func(int) { receiver3 = true },
	}
	broadcast := property.Broadcast[int](receivers, delegation)
	var any int
	broadcast.Change(any)
	if !receiver1 {
		t.Errorf("receiver1: did not get notified")
	}
	if !receiver2 {
		t.Errorf("receiver2: did not get notified")
	}
	if !receiver3 {
		t.Errorf("receiver3: did not get notified")
	}
}

func Test_broadcast_Change_delegates_to_receivers_without_mutation(t *testing.T) {
	table := []string{"go", "golang"}
	delegation := delegate[string]{
		t:      t,
		change: func(string) error { return nil },
	}
	for _, expected := range table {
		receivers := []func(string){
			func(got string) {
				if got != expected {
					t.Errorf("recevier1: expected value is (%v) got (%v)", expected, got)
				}
			},
			func(got string) {
				if got != expected {
					t.Errorf("recevier2: expected value is (%v) got (%v)", expected, got)
				}
			},
		}
		broadcast := property.Broadcast[string](receivers, delegation)
		broadcast.Change(expected)
	}
}

func Test_broadcast_Change_returns_error_without_mutation(t *testing.T) {
	table := []error{
		nil,
		fmt.Errorf("any error"),
	}
	for _, expected := range table {
		delegation := delegate[string]{
			t:      t,
			change: func(string) error { return expected },
		}
		receivers := []func(string){
			func(string) { /* any implementation */ },
		}
		broadcast := property.Broadcast[string](receivers, delegation)
		var any string
		got := broadcast.Change(any)
		if got != expected {
			t.Errorf("expected error is (%v) got (%v)", expected, got)
		}
	}
}

func Test_broadcast_Value_delegates_to_underlying_property(t *testing.T) {
	var invoked bool
	delegation := delegate[int]{
		t: t,
		value: func() (int, error) {
			invoked = true
			var anyNum int
			var anyErr error
			return anyNum, anyErr
		},
	}
	receivers := []func(int){
		func(int) {},
	}
	broadcast := property.Broadcast[int](receivers, delegation)
	broadcast.Value()
	if !invoked {
		t.Error("delegation did not occur to underlying property")
	}
}

func Test_broadcast_Value_returns_value_of_underlying_property(t *testing.T) {
	table := []struct {
		value string
		err   error
	}{
		{"Go", nil},
		{"Golang", fmt.Errorf("any error")},
	}
	for _, data := range table {
		delegation := delegate[string]{
			t: t,
			value: func() (string, error) {
				return data.value, data.err
			},
		}
		receivers := []func(string){
			func(string) {},
		}
		broadcast := property.Broadcast[string](receivers, delegation)
		got, _ := broadcast.Value()
		expected := data.value
		if got != expected {
			t.Errorf("expected value is (%v) got (%v)", expected, got)
		}
	}
}

func Test_broadcast_Value_returns_error_of_underlying_property(t *testing.T) {
	table := []struct {
		value string
		err   error
	}{
		{"Go", nil},
		{"Golang", fmt.Errorf("any error")},
	}
	for _, data := range table {
		delegation := delegate[string]{
			t: t,
			value: func() (string, error) {
				return data.value, data.err
			},
		}
		receivers := []func(string){
			func(string) {},
		}
		broadcast := property.Broadcast[string](receivers, delegation)
		_, got := broadcast.Value()
		expected := data.err
		if got != expected {
			t.Errorf("expected error is (%v) got (%v)", expected, got)
		}
	}
}
