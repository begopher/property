package property

func Guard[T any](cons Constraint[T], property Property[T]) guard[T] {
	if cons == nil {
		panic("property.Guard: cannot be created from nil constraint")
	}
	if property == nil {
		panic("property.Guard: cannot be created from nil property")
	}
	return guard[T]{
		constraint:     cons,
		property: property,
	}
}

type guard[T any] struct {
	constraint     Constraint[T]
	property Property[T]
}

func (g guard[T]) Change(value T) error {
	if err := g.constraint.Evaluate(value); err != nil {
		return err
	}
	return g.property.Change(value)
}

func (g guard[T]) Value() (T, error) {
	return g.property.Value()
}
