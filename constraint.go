package property

type Constraint[T any] interface {
	Evaluate(T) error
}
