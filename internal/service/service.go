package service

import (
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type Repository interface {
	Save(value model.Metrics)
	Get(id string, mType string) *model.Metrics
	GetAll() []model.Metrics
}

type DBClient interface {
	Ping() error
}

type Service struct {
	storage Repository
	client  DBClient
}

func NewService(storage Repository, client DBClient) *Service {
	service := &Service{
		storage: storage,
		client:  client,
	}

	return service
}

func (service *Service) DBPing() error {
	return service.client.Ping()
}

func (service *Service) Save(metric model.Metrics) {
	storage := service.storage
	existedMetric := storage.Get(metric.ID, metric.MType)

	if metric.MType == model.Counter {
		if existedMetric != nil {
			delta := *existedMetric.Delta + int64(*metric.Delta)
			existedMetric.Delta = &delta
			storage.Save(*existedMetric)
		} else {
			storage.Save(metric)
		}
	} else {
		storage.Save(metric)
	}
}

func (service *Service) Get(id string, mType string) *model.Metrics {
	storage := service.storage
	metric := storage.Get(id, mType)
	return metric
}

func (service *Service) GetAll() []model.Metrics {
	storage := service.storage
	return storage.GetAll()
}
