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

	metric := container.service.Get(metricName, metricType)

	if metric == nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	switch metric.MType {
	case model.Counter:
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%d", *metric.Delta)
	case model.Gauge:
		res.WriteHeader(http.StatusOK)
		fmt.Fprintf(res, "%g", *metric.Value)
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

	metricValue := container.service.Get(metric.ID, metric.MType)

	if metricValue == nil {
		logger.Error("Metric not found")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	returnBody, err := json.Marshal(metricValue)

	if err == nil {
		res.WriteHeader(http.StatusOK)
		res.Write(returnBody)
		return
	}

	logger.Error("Failed to write body")
	logger.Error(err.Error())
	res.WriteHeader(http.StatusInternalServerError)
}
