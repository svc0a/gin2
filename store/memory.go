package store

import (
	"github.com/svc0a/worker/syncx"
)

type MemoryStore[T any] struct {
	data syncx.Map[T]
}

func NewMemoryStore[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{data: syncx.Define[T]()}
}

func (m *MemoryStore[T]) Store(key string, value T) {
	m.data.Store(key, value)
}

func (m *MemoryStore[T]) Load(key string) (*T, error) {
	data, err := m.data.Load(key)
	if err != nil {
		return nil, err
	}
	return data, nil
}
