package middleware

import (
	"github.com/svc0a/worker/syncx"
)

type MemoryStore struct {
	data syncx.Map[[]byte]
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: syncx.Define[[]byte]()}
}

func (m *MemoryStore) Store(key string, value []byte) {
	m.data.Store(key, value)
}

func (m *MemoryStore) Load(key string) ([]byte, error) {
	data, err := m.data.Load(key)
	if err != nil {
		return nil, err
	}
	return *data, nil
}
