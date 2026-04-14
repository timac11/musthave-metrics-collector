package handler

import (
	"encoding/json"
	"net/http"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

func (container *ApplicationAPIContainer) UpdateMetrics(res http.ResponseWriter, req *http.Request) {
	metrics, validationRes := parseMetrics(req)

	if validationRes != nil {
		logger.Error("Validation error", "message", validationRes.Message)
		http.Error(res, validationRes.Message, validationRes.Code)
		return
	}

	err := container.service.SaveAll(req.Context(), metrics)

	if err != nil {
		logger.Error("Internal server error", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func parseMetrics(req *http.Request) ([]*model.Metrics, *model.ValidationErr) {
	var metrics []*model.Metrics

	defer req.Body.Close()

	err := json.NewDecoder(req.Body).Decode(&metrics)

	if err != nil {
		return nil, &model.ValidationErr{Message: "Invalid metric", Code: http.StatusBadRequest}
	}

	for _, metric := range metrics {
		if metric.MType != model.Counter && metric.MType != model.Gauge {
			return nil, &model.ValidationErr{Message: "Invalid metric type", Code: http.StatusBadRequest}
		}

		if metric.MType == model.Counter && metric.Delta == nil {
			return nil, &model.ValidationErr{Message: "Invalid counter value", Code: http.StatusBadRequest}
		}

		if metric.MType == model.Gauge && metric.Value == nil {
			return nil, &model.ValidationErr{Message: "Invalid gauge value", Code: http.StatusBadRequest}
		}
	}

	return metrics, nil
}
