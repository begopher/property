package test

import (
	"testing"
)

type delegate[T any] struct {
	t      *testing.T
	change func(T) error
	value  func() (T, error)
}

func (d delegate[T]) Change(value T) error {
	if d.change == nil {
		d.t.Fatalf("underlying property.Change(%v) method has invoked", value)
		return nil
	}
	return d.change(value)
}

func (d delegate[T]) Value() (T, error) {
	if d.value == nil {
		d.t.Fatalf("underlying property.Value() method has invoked")
		var value T
		return value, nil
	}
	return d.value()
}
