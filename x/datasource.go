package x

type Datasource[T any] interface {
	Change(T) error
	Value() (T, error)
}
