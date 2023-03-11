package property

type Rule[T any] interface {
	Evaluate(T) error
}
