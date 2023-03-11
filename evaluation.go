package property

func Evaluation[T any](rule Rule[T], property Property[T]) evaluation[T] {
	if rule == nil {
		panic("property.Evaluation: cannot be created from nil rule")
	}
	if property == nil {
		panic("property.Evaluation: cannot be created from nil property")
	}
	return evaluation[T]{
		rule:     rule,
		property: property,
	}
}

type evaluation[T any] struct {
	rule     Rule[T]
	property Property[T]
}

func (c evaluation[T]) Change(value T) error {
	if err := c.rule.Evaluate(value); err != nil {
		return err
	}
	return c.property.Change(value)
}

func (c evaluation[T]) Value() (T, error) {
	return c.property.Value()
}
