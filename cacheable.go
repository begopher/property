package property

func Cacheable[T any](property Property[T]) *cacheable[T] {
	if property == nil {
		panic("property.Cacheable: cannot be created from nil property")
	}
	return &cacheable[T]{
		property: property,
	}
}

type cacheable[T any] struct {
	value    *T
	property Property[T]
}

func (c *cacheable[T]) Change(value T) error {
	err := c.property.Change(value)
	if err == nil {
		c.value = &value
	}
	return err
}

func (c *cacheable[T]) Value() (T, error) {
	if c.value != nil {
		return *c.value, nil
	}
	value, err := c.property.Value()
	if err == nil {
		c.value = &value
	}
	return value, err
}
