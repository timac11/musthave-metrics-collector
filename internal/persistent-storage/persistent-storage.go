package persistentstorage

import (
	"encoding/json"
	"os"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type FileStorage struct {
	path string
}

func NewPersistentStorage(path string) *FileStorage {
	fileStorage := FileStorage{path: path}
	return &fileStorage
}

func (fs *FileStorage) Store(value map[string]*model.Metrics) error {
	logger.Info("Start backup data to file")

	file, err := os.OpenFile(fs.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		logger.Error("Failed to open backup file", err.Error())
		return err
	}
	defer file.Close()

	data, err := json.Marshal(&value)

	if err != nil {
		logger.Error("Failed to marshal backup file", err.Error())
		return err
	}

	file.Write(data)

	return nil
}

func (fs *FileStorage) Restore() (map[string]*model.Metrics, error) {
	data, err := os.ReadFile(fs.path)

	if err != nil {
		logger.Error("Failed to read backup file", err.Error())
		return nil, err
	}

	var value map[string]*model.Metrics
	err = json.Unmarshal(data, &value)

	if err != nil {
		logger.Error("Failed to unmarshal backup file", err.Error())
		return nil, err
	}

	return value, nil
}
