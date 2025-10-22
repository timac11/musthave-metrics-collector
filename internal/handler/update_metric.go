package handler

import (
	"fmt"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"net/http"
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

func parseMetricParams(req *http.Request) (*model.MetricInfo, *model.ValidationErr) {
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
		var value float64
		_, err := fmt.Sscanf(metricValue, "%f", &value)
		if err != nil {
			return nil, &model.ValidationErr{Message: "Invalid gauge value", Code: http.StatusBadRequest}
		}

		return &model.MetricInfo{Name: metricName, MType: metricType, Value: &value}, nil
	case model.Counter:
		var value int64
		_, err := fmt.Sscanf(metricValue, "%d", &value)
		if err != nil {
			return nil, &model.ValidationErr{Message: "Invalid counter value", Code: http.StatusBadRequest}
		}
		return &model.MetricInfo{Name: metricName, MType: metricType, Delta: &value}, nil
	}

	return nil, &model.ValidationErr{Message: "Invalid metric type", Code: http.StatusBadRequest}
}
