package repository

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type MemStorage struct {
	storage    map[string]model.Metrics
	mu         *sync.Mutex
	backupPath string
}

func NewMemStorage(backupPath string) *MemStorage {
	mu := sync.Mutex{}
	ms := &MemStorage{
		storage:    make(map[string]model.Metrics),
		mu:         &mu,
		backupPath: backupPath,
	}

	return ms
}

func (ms *MemStorage) Save(metric model.Metrics) {
	// save if does not exist and rewrite if exist
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.storage[buildMetricHash(metric)] = metric
	ms.backup()
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

func (ms *MemStorage) Restore() error {
	data, err := os.ReadFile(ms.backupPath)
	if err != nil {
		logger.Error("Failed to read backup file", err.Error())
		return err
	}

	var memsMap map[string]model.Metrics
	err = json.Unmarshal(data, &memsMap)

	if err != nil {
		logger.Error("Failed to unmarshal backup file", err.Error())
		return err
	}

	ms.storage = memsMap

	return nil
}

func (ms *MemStorage) backup() error {
	logger.Info("Start backup data to file")

	file, err := os.OpenFile(ms.backupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		logger.Error("Failed to open backup file", err.Error())
		return err
	}
	defer file.Close()

	data, err := json.Marshal(&ms.storage)

	if err != nil {
		logger.Error("Failed to marshal backup file", err.Error())
		return err
	}

	file.Write(data)

	return nil
}
