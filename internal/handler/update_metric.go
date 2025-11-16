package handler

import (
	"encoding/json"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func (container *ApplicationAPIContainer) UpdateMetric(res http.ResponseWriter, req *http.Request) {
	metric, validationRes := parseMetricParams(req)

	if validationRes != nil {
		http.Error(res, validationRes.Message, validationRes.Code)
		return
	}

	container.service.Save(*metric)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func (container *ApplicationAPIContainer) UpdateMetricV2(res http.ResponseWriter, req *http.Request) {
	metric, validationRes := parseMetricParamsV2(req)

	if validationRes != nil {
		http.Error(res, validationRes.Message, validationRes.Code)
		return
	}

	container.service.Save(*metric)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func parseMetricParamsV2(req *http.Request) (*model.Metrics, *model.ValidationErr) {
	var metric model.Metrics

	err := json.NewDecoder(req.Body).Decode(&metric)
	if err != nil {
		return nil, &model.ValidationErr{Message: "Invalid metric", Code: http.StatusBadRequest}
	}

	if metric.MType != model.Counter && metric.MType != model.Gauge {
		return nil, &model.ValidationErr{Message: "Invalid metric type", Code: http.StatusBadRequest}
	}

	return &metric, nil
}

func parseMetricParams(req *http.Request) (*model.Metrics, *model.ValidationErr) {
	path := req.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) != 4 {
		return nil, &model.ValidationErr{Message: "Invalid URL format", Code: http.StatusNotFound}
	}

	metricType := parts[1]
	metricName := parts[2]
	metricValue := parts[3]

	switch metricType {
	case model.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return nil, &model.ValidationErr{Message: "Invalid gauge value", Code: http.StatusBadRequest}
		}

		return &model.Metrics{ID: metricName, MType: metricType, Value: &value}, nil
	case model.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return nil, &model.ValidationErr{Message: "Invalid counter value", Code: http.StatusBadRequest}
		}
		return &model.Metrics{ID: metricName, MType: metricType, Delta: &value}, nil
	}

	return nil, &model.ValidationErr{Message: "Invalid metric type", Code: http.StatusBadRequest}
}
