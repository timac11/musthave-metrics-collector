package repository

import (
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"sync"
)

type MemStorage struct {
	storage map[string]model.Metrics
	mu      *sync.Mutex
}

func (ms *MemStorage) Save(value model.Metrics) {
	// save if does not exist and rewrite if exist
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.storage[value.ID] = value
}

func (ms *MemStorage) Get(key string) *model.Metrics {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	val, ok := ms.storage[key]
	if ok {
		return &val
	}
	return nil
}

func (ms *MemStorage) GetAll() []model.Metrics {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	metrics := make([]model.Metrics, 0, len(ms.storage))

	for _, v := range ms.storage {
		metrics = append(metrics, v)
	}

	return metrics
}

func NewMemStorage() *MemStorage {
	mu := sync.Mutex{}
	ms := &MemStorage{
		storage: make(map[string]model.Metrics),
		mu:      &mu,
	}

	return ms
}
