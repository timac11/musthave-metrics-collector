package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"net/http"
)

func (container *ApplicationAPIContainer) GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	metric := container.service.Get(metricType, metricName)

	switch v := metric.(type) {
	case int64:
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%d", v)
	case float64:
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%g", v)
	default:
		res.WriteHeader(http.StatusNotFound)
	}
}

func (container *ApplicationAPIContainer) GetFullMetricInfo(res http.ResponseWriter, req *http.Request) {
	var metric model.Metrics
	res.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(req.Body).Decode(&metric)
	if err != nil {
		logger.Error("Failed to decode body")
		logger.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	metricValue := container.service.Get(metric.MType, metric.ID)

	switch v := metricValue.(type) {
	case int64:
		metric.Delta = &v
	case float64:
		metric.Value = &v
	default:
		logger.Info("Metric not found")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	valueMetric, err := json.Marshal(metric)
	if err == nil {
		res.Write(valueMetric)
		return
	}
	logger.Error("Failed to write body")
	logger.Error(err.Error())
	res.WriteHeader(http.StatusInternalServerError)
}
