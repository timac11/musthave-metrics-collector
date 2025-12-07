package service

import (
	"context"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type Repository interface {
	Save(ctx context.Context, value model.Metrics) error
	Get(ctx context.Context, id string, mType string) (*model.Metrics, error)
	GetAll(ctx context.Context) ([]model.Metrics, error)
	Ping(ctx context.Context) error
}

type Service struct {
	storage Repository
}

func NewService(storage Repository) *Service {
	service := &Service{
		storage: storage,
	}

	return service
}

func (service *Service) DBPing() error {
	return service.storage.Ping(context.Background())
}

func (service *Service) Save(metric model.Metrics) {
	storage := service.storage

	if metric.MType == model.Counter {
		existedMetric, err := storage.Get(context.Background(), metric.ID, metric.MType)

		if err == nil && existedMetric != nil {
			delta := *existedMetric.Delta + int64(*metric.Delta)
			existedMetric.Delta = &delta

			logger.Info("updated metric", existedMetric.ID, *existedMetric.Delta)

			storage.Save(context.Background(), *existedMetric)
		} else {
			storage.Save(context.Background(), metric)
		}
	} else {
		storage.Save(context.Background(), metric)
	}
}

func (service *Service) Get(id string, mType string) *model.Metrics {
	storage := service.storage
	metric, err := storage.Get(context.Background(), id, mType)

	if err == nil {
		return nil
	}

	return metric
}

func (service *Service) GetAll() []model.Metrics {
	storage := service.storage
	metrics, err := storage.GetAll(context.Background())
	if err != nil {
		return []model.Metrics{}
	}

	return metrics
}
