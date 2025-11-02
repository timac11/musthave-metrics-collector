package service

import (
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type Repository interface {
	Save(value model.Metrics)
	Get(key string) *model.Metrics
	GetAll() []model.Metrics
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

func (service *Service) Save(metric model.MetricInfo) {
	storage := service.storage
	metricID := buildMetricID(metric.MType, metric.Name)
	existedMetric := storage.Get(metricID)

	if metric.MType == model.Counter {
		if existedMetric != nil {
			delta := *existedMetric.Delta + int64(*metric.Delta)
			updatedMetric := model.Metrics{
				ID:    existedMetric.ID,
				Hash:  existedMetric.Hash,
				Delta: &delta,
				MType: existedMetric.MType,
			}

			storage.Save(updatedMetric)
		} else {
			delta := int64(*metric.Delta)
			metric := model.Metrics{
				ID:    metricID,
				Hash:  metricID,
				Delta: &delta,
				MType: model.Counter,
			}

			storage.Save(metric)
		}
	} else {
		metric := model.Metrics{
			ID:    metricID,
			Hash:  metricID,
			Value: metric.Value,
			MType: model.Gauge,
		}

		storage.Save(metric)
	}
}

func (service *Service) Get(metricType string, metricName string) interface{} {
	storage := service.storage
	metric := storage.Get(buildMetricID(metricType, metricName))

	if metric == nil {
		return nil
	}

	if metric.MType == model.Counter {
		return *metric.Delta
	}
	return *metric.Value
}

func (service *Service) GetAll() []model.Metrics {
	storage := service.storage
	return storage.GetAll()
}

func buildMetricID(metricType string, metricName string) string {
	return metricType + "-" + metricName
}
