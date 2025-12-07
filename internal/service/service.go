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

func (service *Service) Save(metric model.Metrics) error {
	storage := service.storage
	var err error

	if metric.MType == model.Counter {
		var existedMetric *model.Metrics
		existedMetric, err = storage.Get(context.Background(), metric.ID, metric.MType)

		if err != nil {
			return nil
		}

		if existedMetric != nil {
			delta := *existedMetric.Delta + int64(*metric.Delta)
			existedMetric.Delta = &delta

			logger.Info("updated metric", existedMetric.ID, *existedMetric.Delta)

			err = storage.Save(context.Background(), *existedMetric)
		} else {
			err = storage.Save(context.Background(), metric)
		}
	} else {
		storage.Save(context.Background(), metric)
	}

	return err
}

func (service *Service) Get(id string, mType string) (*model.Metrics, error) {
	storage := service.storage
	metric, err := storage.Get(context.Background(), id, mType)

	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (service *Service) GetAll() ([]model.Metrics, error) {
	storage := service.storage
	metrics, err := storage.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
