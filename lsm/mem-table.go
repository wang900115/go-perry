package main

import "sync"

type Value = []byte

type MemTable struct {
	mu   sync.RWMutex
	data map[string]Value
}

func NewMemTable() *MemTable {
	return &MemTable{
		data: make(map[string]Value),
	}
}

func (m *MemTable) Put(key string, value Value) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make(Value, len(value))
	copy(cp, value)
	m.data[key] = cp
}

func (m *MemTable) Get(key string) (Value, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[key]
	if !ok {
		return nil, false
	}
	cp := make(Value, len(value))
	copy(cp, value)
	return cp, true
}

func (m *MemTable) Flush() map[string]Value {
	m.mu.Lock()
	defer m.mu.Unlock()
	flushedData := m.data
	m.data = make(map[string]Value)
	return flushedData
}
