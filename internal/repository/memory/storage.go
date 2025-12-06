package memorystorage

import (
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"sync"
)

type PersistentStorage interface {
	Store(value map[string]model.Metrics) error
	Restore() (map[string]model.Metrics, error)
}

type MemStorage struct {
	storage           map[string]model.Metrics
	mu                *sync.Mutex
	persistentStorage PersistentStorage
}

func NewMemStorage(ps PersistentStorage, restore bool) *MemStorage {
	mu := sync.Mutex{}
	storage := make(map[string]model.Metrics)

	if restore {
		metrics, err := ps.Restore()
		if err != nil {
			logger.Error("Failed restore metrics")
			logger.Error(err.Error())
		} else {
			storage = metrics
		}
	}

	ms := &MemStorage{
		storage:           storage,
		mu:                &mu,
		persistentStorage: ps,
	}

	return ms
}

func (ms *MemStorage) Save(metric model.Metrics) {
	// save if does not exist and rewrite if exist
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.storage[buildMetricHash(metric)] = metric
	ms.persistentStorage.Store(ms.storage)
}

func buildMetricHash(metric model.Metrics) string {
	return metric.ID + "-" + metric.MType
}

func (ms *MemStorage) Get(id string, mType string) *model.Metrics {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	val, ok := ms.storage[id+"-"+mType]
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
