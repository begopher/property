package property

func Inequality[T comparable](property Property[T]) inequality[T] {
	return inequality[T]{property}
}

type inequality[T comparable] struct {
	property Property[T]
}

func (i inequality[T]) Change(value T) error {
	old, err := i.Value()
	if err != nil {
		return err
	}
	if old == value {
		return nil
	}
	return i.property.Change(value)
}

func (i inequality[T]) Value() (T, error) {
	return i.property.Value()
}

