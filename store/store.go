package store

type Store[T any] interface {
	Store(key string, value *T)
	Load(key string) (*T, error)
}
