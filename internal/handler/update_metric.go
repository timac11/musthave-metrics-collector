package handler

import (
	"fmt"
	"net/http"
	"strings"

	logger "github.com/timac11/musthave-metrics-collector/internal/logger"
	model "github.com/timac11/musthave-metrics-collector/internal/model"
)

func UpdateMetric(res http.ResponseWriter, req *http.Request) {
	_, validationRes := parseMetricParams(req.URL.Path)

	if validationRes != nil {
		http.Error(res, validationRes.Message, validationRes.Code)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func parseMetricParams(path string) (*model.Metrics, *model.ValidationErr) {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) != 4 {
		return nil, &model.ValidationErr{Message: "Invalid URL format", Code: http.StatusNotFound}
	}

	metricType := parts[1]
	metricName := parts[2]
	metricValue := parts[3]

	// TODO: add to switch case
	if metricType != model.Gauge && metricType != model.Counter {
		return nil, &model.ValidationErr{Message: "Invalid metric type", Code: http.StatusBadRequest}
	}

	switch metricType {
	case model.Gauge:
		var value float64
		_, err := fmt.Sscanf(metricValue, "%f", &value)
		if err != nil {
			return nil, &model.ValidationErr{Message: "Invalid gauge value", Code: http.StatusBadRequest}
		}
	case model.Counter:
		var value int64
		_, err := fmt.Sscanf(metricValue, "%d", &value)
		if err != nil {
			return nil, &model.ValidationErr{Message: "Invalid counter value", Code: http.StatusBadRequest}
		}
	}

	var floatValue float64
	fmt.Sscanf(metricValue, "%f", &floatValue)

	logger.Log("metric value")
	logger.Log(floatValue)

	return &model.Metrics{ID: metricName, MType: metricType, Value: &floatValue}, nil
}
