package property

// Cache implements Property[T any] interface which provides the
// ability to save value in memory.
// Cache represents eager loading since value must be provided
// at construction time.
//
// messages:
//   - Change updates underlying property and in memory value if no error occur.
//   - Value  returns in-memory cached value and nil.
//
// # Panic when property argument is nil
//
// See Cacheable if lazy loading is needed.
func Cache[T any](value T, property Property[T]) *cache[T] {
	if property == nil {
		panic("property.Cache: cannot be created from nil property")
	}
	return &cache[T]{
		value:    value,
		property: property,
	}
}

// Cache implements Property[T any] interface
type cache[T any] struct {
	value    T
	property Property[T]
}

// Change message delegates to underlying property to update itself
// when no error occur, in-memory cached value will be updated.
//
// Error of underlying property is returned.
func (c *cache[T]) Change(value T) error {
	err := c.property.Change(value)
	if err == nil {
		c.value = value
	}
	return err
}

// Value message returns in-memory cached value and nil.
// No delegation occurs to underlying property.
func (c *cache[T]) Value() (T, error) {
	return c.value, nil
}
