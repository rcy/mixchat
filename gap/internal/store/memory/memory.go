package memory

import (
	"context"
	"fmt"
)

type Memory struct {
	cache map[string][]byte
}

func MustInit() *Memory {
	fmt.Printf("Initialized memory based storage\n")
	return &Memory{cache: map[string][]byte{}}
}

func (m *Memory) Put(ctx context.Context, key string, data []byte) error {
	m.cache[key] = data
	return nil
}

func (m *Memory) Get(ctx context.Context, key string) ([]byte, error) {
	data, ok := m.cache[key]
	if ok {
		return data, nil
	} else {
		return nil, fmt.Errorf("not found")
	}
}
