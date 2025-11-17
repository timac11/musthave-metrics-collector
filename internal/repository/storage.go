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

func (ms *MemStorage) Save(value model.Metrics) {
	// save if does not exist and rewrite if exist
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.storage[value.ID+"-"+value.MType] = value
	ms.backup()
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
	logger.Info("Start restore data")

	file, err := os.OpenFile(ms.backupPath, os.O_RDONLY|os.O_CREATE, 0755)

	if err != nil {
		logger.Error("Failed to open backup file", err)
		return err
	}
	defer file.Close()

	var data []byte

	file.Read(data)

	memsMap := map[string]model.Metrics{}
	err = json.Unmarshal(data, &memsMap)

	if err != nil {
		logger.Error("Failed to unmarshal backup file", err)
		return err
	}

	ms.storage = memsMap

	return nil
}

func (ms *MemStorage) backup() error {
	logger.Info("Start backup data to file")

	file, err := os.OpenFile(ms.backupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		logger.Error("Failed to open backup file", err)
		return err
	}
	defer file.Close()

	data, marshalErr := json.Marshal(&ms.storage)

	if marshalErr != nil {
		logger.Error("Failed to marshal backup file", marshalErr)
		return marshalErr
	}

	file.Write(data)

	return nil
}
