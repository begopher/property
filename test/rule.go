package test

type rule[T any] struct {
	evaluate func(T) error
}

func (r rule[T]) Evaluate(value T) error {
	return r.evaluate(value)
}
