package memorystorage

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	persistentstorage "github.com/timac11/musthave-metrics-collector/internal/persistent-storage"
)

func TestSaveMetricToStorage(t *testing.T) {
	firstValue := float64(123.456)
	secondValue := int64(567)

	tests := []model.Metrics{
		{
			ID:    "first",
			MType: model.Gauge,
			Value: &firstValue,
		},
		{
			ID:    "second",
			MType: model.Counter,
			Delta: &secondValue,
		},
	}

	for _, test := range tests {
		t.Run(test.ID, func(t *testing.T) {
			dbFilePath := "/tmp/db" + test.ID + ".json"

			persistentStorage := persistentstorage.NewPersistentStorage(dbFilePath)
			storage := NewMemStorage(persistentStorage, false)

			storage.Save(context.Background(), &test)
			savedMetric, _ := storage.Get(context.Background(), test.ID, test.MType)

			assert.Equal(t, savedMetric.ID, test.ID)
			assert.Equal(t, savedMetric.MType, test.MType)
			assert.Equal(t, savedMetric.Value, test.Value)
			assert.Equal(t, savedMetric.Delta, test.Delta)

			os.Remove(dbFilePath)
		})
	}
}

func TestUpdatMetricInStorage(t *testing.T) {
	firstValue := float64(123.456)
	secondValue := float64(567)

	metric := &model.Metrics{
		ID:    "first",
		MType: model.Gauge,
		Value: &firstValue,
	}

	dbFilePath := "/tmp/db_test_update.json"
	defer os.Remove(dbFilePath)

	persistentStorage := persistentstorage.NewPersistentStorage(dbFilePath)
	storage := NewMemStorage(persistentStorage, false)

	storage.Save(context.Background(), metric)
	savedMetric, _ := storage.Get(context.Background(), metric.ID, metric.MType)

	assert.Equal(t, savedMetric.ID, metric.ID)
	assert.Equal(t, savedMetric.MType, metric.MType)
	assert.Equal(t, savedMetric.Value, metric.Value)
	assert.Equal(t, savedMetric.Delta, metric.Delta)

	allMetrics, _ := storage.GetAll(context.Background())

	assert.Equal(t, len(allMetrics), 1)

	metric.Value = &secondValue
	storage.Save(context.Background(), metric)
	savedMetric, _ = storage.Get(context.Background(), metric.ID, metric.MType)

	assert.Equal(t, *savedMetric.Value, secondValue)

	allMetrics, _ = storage.GetAll(context.Background())

	assert.Equal(t, len(allMetrics), 1)
}

func TestBackupInStorage(t *testing.T) {
	dbFilePath := "/tmp/db_test_backup.json"
	persistentStorage := persistentstorage.NewPersistentStorage(dbFilePath)
	storage := NewMemStorage(persistentStorage, false)

	defer os.Remove(dbFilePath)

	value := float64(123.456)

	metric := &model.Metrics{
		ID:    "first",
		MType: model.Gauge,
		Value: &value,
	}

	storage.Save(context.Background(), metric)

	data, err := os.ReadFile(dbFilePath)
	assert.Nil(t, err)

	var memsMap map[string]model.Metrics
	err = json.Unmarshal(data, &memsMap)
	assert.Nil(t, err)

	backupedMetric := memsMap[buildMetricHash(metric)]

	assert.NotNil(t, backupedMetric)

	assert.Equal(t, metric.ID, backupedMetric.ID)
	assert.Equal(t, metric.MType, backupedMetric.MType)
	assert.Equal(t, metric.Value, backupedMetric.Value)
	assert.Equal(t, metric.Delta, backupedMetric.Delta)
}

func TestRestoreInStorage(t *testing.T) {
	value := float64(123.456)
	metric := &model.Metrics{
		ID:    "first",
		MType: model.Gauge,
		Value: &value,
	}

	dbFilePath := "/tmp/db_test_restore.json"
	defer os.Remove(dbFilePath)
	file, err := os.Create(dbFilePath)
	assert.Nil(t, err)

	// create mems map and save to file
	memsMap := make(map[string]*model.Metrics)

	memsMap[buildMetricHash(metric)] = metric

	data, err := json.Marshal(&memsMap)
	assert.Nil(t, err)

	file.Write(data)
	file.Close()

	// restore metrics in storage
	persistentStorage := persistentstorage.NewPersistentStorage(dbFilePath)
	storage := NewMemStorage(persistentStorage, true)

	restoredMetric, _ := storage.Get(context.Background(), metric.ID, metric.MType)

	assert.Equal(t, metric.ID, restoredMetric.ID)
	assert.Equal(t, metric.MType, restoredMetric.MType)
	assert.Equal(t, metric.Value, restoredMetric.Value)
	assert.Equal(t, metric.Delta, restoredMetric.Delta)
}
