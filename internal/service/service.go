package service

import (
	model "github.com/timac11/musthave-metrics-collector/internal/model"
	repository "github.com/timac11/musthave-metrics-collector/internal/repository"
)

type Service struct {
	storage repository.Repository
}

func NewService(s repository.Repository) *Service {
	service := &Service{
		storage: s,
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

func (service *Service) Get(metricType string, metricName string) *model.Metrics {
	storage := service.storage
	return storage.Get(buildMetricID(metricType, metricName))
}

func (service *Service) GetAll() []model.Metrics {
	storage := service.storage
	return storage.GetAll()
}

func buildMetricID(metricType string, metricName string) string {
	return metricType + "-" + metricName
}
