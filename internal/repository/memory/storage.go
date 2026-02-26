package memorystorage

import (
	"context"
	"errors"
	"sync"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type PersistentStorage interface {
	Store(value map[string]*model.Metrics) error
	Restore() (map[string]*model.Metrics, error)
}

type MemStorage struct {
	storage           map[string]*model.Metrics
	mu                *sync.Mutex
	persistentStorage PersistentStorage
}

func NewMemStorage(ps PersistentStorage, restore bool) *MemStorage {
	mu := sync.Mutex{}
	storage := make(map[string]*model.Metrics)

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

func (ms *MemStorage) Save(ctx context.Context, metric *model.Metrics) error {
	existedMetric, _ := ms.Get(ctx, metric.ID, metric.MType)

	if existedMetric != nil {
		if existedMetric.MType == model.Counter && existedMetric.Delta != nil {
			newDelta := *metric.Delta + *existedMetric.Delta
			metric.Delta = &newDelta
		}
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.storage[buildMetricHash(metric)] = metric
	ms.persistentStorage.Store(ms.storage)

	return nil
}

func (ms *MemStorage) SaveAll(ctx context.Context, metrics []*model.Metrics) error {
	for _, metric := range metrics {
		err := ms.Save(ctx, metric)

		if err != nil {
			return err
		}
	}

	return nil
}

func (ms *MemStorage) Ping(_ context.Context) error {
	return errors.New("connection was not established")
}

func (ms *MemStorage) Get(_ context.Context, id string, mType string) (*model.Metrics, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	val, ok := ms.storage[id+"-"+mType]
	if ok {
		return val, nil
	}
	return nil, errors.New("failed to get metric")
}

func (ms *MemStorage) GetAll(_ context.Context) ([]*model.Metrics, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	metrics := make([]*model.Metrics, 0, len(ms.storage))

	for _, v := range ms.storage {
		metrics = append(metrics, v)
	}

	return metrics, nil
}

func buildMetricHash(metric *model.Metrics) string {
	return metric.ID + "-" + metric.MType
}
