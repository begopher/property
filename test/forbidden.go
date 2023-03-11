package test

import (
	"testing"
)

type forbidden[T any] struct {
	t *testing.T
}

func (f forbidden[T]) Change(value T) error {
	f.t.Errorf("forbidden.Change(%v) has been invoked", value)
	return nil
}

func (f forbidden[T]) Value() (T, error) {
	f.t.Error("forbidden.Value() has been invoked")
	var value T
	return value, nil
}
