package repository

import (
	model "github.com/timac11/musthave-metrics-collector/internal/model"
	"maps"
)

type Repository interface {
	Save(value model.Metrics)
	Get(key string) *model.Metrics
	GetAll() []model.Metrics
}

type MemStorage struct {
	storage map[string]model.Metrics
}

func (ms *MemStorage) Save(value model.Metrics) {
	// save if does not exist and rewrite if exist
	ms.storage[value.ID] = value
}

func (ms *MemStorage) Get(key string) *model.Metrics {
	val, ok := ms.storage[key]
	if ok {
		return &val
	}
	return nil
}

func (ms *MemStorage) GetAll() []model.Metrics {
	metricsMap := maps.Clone(ms.storage)
	metrics := make([]model.Metrics, 0, len(metricsMap))

	for _, v := range metricsMap {
		metrics = append(metrics, v)
	}

	return metrics
}

func NewMemStorage() *MemStorage {
	ms := &MemStorage{
		storage: make(map[string]model.Metrics),
	}

	return ms
}
