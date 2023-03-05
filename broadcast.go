package property

// Broadcast returns implementation of Property[T any] interface.
// it is responsible for notifing all revicers synchronously when
// value of underlying property get updated.
//
// Broadcast is useful when a property represents a natural primary key,
// which has to change and any other properties that depends upon the old
// primary key must get updated.
//
// receivers is a slice of functions so client can choose appropriate message name
// panic when:
//   - property is nil
//   - recivers is nil.
//   - recivers has nil function (default panic message).
func Broadcast[T any](receivers []func(T), property Property[T]) *broadcast[T] {
	if len(receivers) == 0 {
		panic("property.broadcast: cannot be created with zero receivers")
	}
	if property == nil {
		panic("property.broadcast: cannot be created from nil property")
	}
	return &broadcast[T]{
		receivers: receivers,
		property:  property,
	}
}

type broadcast[T any] struct {
	receivers []func(T)
	property  Property[T]
}

// Change message delegates to underlying property to update itself
// when no error occur, receivers get notified with the same value
//
// Error of underlying property is returned.
func (b broadcast[T]) Change(value T) error {
	if err := b.property.Change(value); err != nil {
		return err
	}
	for _, fx := range b.receivers {
		fx(value)
	}
	return nil
}

// Value message returns actual value by delegation to underlying property.
//
// error of underlying of underlying property is returned.
func (b broadcast[T]) Value() (T, error) {
	return b.property.Value()
}
