package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

// GetMetric godoc
// @Summary      Return metric info
// @Description  Get metric by metricType and metricName
// @Tags         Metrics
// @Accept       json
// @Produce      json
// @Param        metricType   path    string  true  "Metric type"
// @Param        metricName   path    string  true  "Metric name"
// @Success      200  {object}  model.Metrics
// @Router       /value/{metricType}/{metricName} [get]
func (container *ApplicationAPIContainer) GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	metric, err := container.service.Get(req.Context(), metricName, metricType)

	if err != nil || metric == nil {
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

// GetFullMetricInfo godoc
// @Summary      Return metric info
// @Description  Get metric by ID and type
// @Tags         Metrics
// @Accept       json
// @Produce      json
// @Param        body  body  model.Metrics  true  "Metric type"
// @Success      200  {object}  model.Metrics
// @Router       /value [post]
func (container *ApplicationAPIContainer) GetFullMetricInfo(res http.ResponseWriter, req *http.Request) {
	var metric model.Metrics

	err := json.NewDecoder(req.Body).Decode(&metric)
	if err != nil {
		logger.Error("Failed to decode body")
		logger.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("Get metric params", metric.ID, metric.MType)

	metricValue, err := container.service.Get(req.Context(), metric.ID, metric.MType)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	returnBody, err := json.Marshal(metricValue)

	if err == nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(returnBody)
		return
	}

	logger.Error("Failed to write body")
	logger.Error(err.Error())
	res.WriteHeader(http.StatusInternalServerError)
}
