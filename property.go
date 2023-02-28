package property

type Property[T any] interface {
	Change(T) error
	Value() (T, error)
}
