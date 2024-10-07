package idempotency

import (
	"github.com/svc0a/worker/syncx"
)

type MemoryStore struct {
	data syncx.Map[*Response]
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: syncx.Define[*Response]()}
}

func (m *MemoryStore) Store(key string, value *Response) {
	m.data.Store(key, value)
}

func (m *MemoryStore) Load(key string) (*Response, error) {
	data, err := m.data.Load(key)
	if err != nil {
		return nil, err
	}
	return *data, nil
}
