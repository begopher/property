package property

func LazyCache[T any](property Property[T]) *lazyCache[T] {
	if property == nil {
		panic("property.LazyCache: cannot be created from nil property")
	}
	return &lazyCache[T]{
		property: property,
	}
}

type lazyCache[T any] struct {
	value    *T
	property Property[T]
}

func (c *lazyCache[T]) Change(value T) error {
	err := c.property.Change(value)
	if err == nil {
		c.value = &value
	}
	return err
}

func (c *lazyCache[T]) Value() (T, error) {
	if c.value != nil {
		return *c.value, nil
	}
	value, err := c.property.Value()
	if err == nil {
		c.value = &value
	}
	return value, err
}
