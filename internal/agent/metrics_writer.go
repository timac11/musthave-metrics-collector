package agent

import (
	"fmt"
	"net/http"
	"strconv"

	logger "github.com/timac11/musthave-metrics-collector/internal/logger"
	model "github.com/timac11/musthave-metrics-collector/internal/model"
)

var baseUrl = "http://localhost:8080"

func Write(metrics []model.Metrics) {
	for _, metric := range metrics {
		writeMetric(metric)
	}
}

func writeMetric(metric model.Metrics) {
	url := fmt.Sprintf("%s/%s/%s/%s", baseUrl, metric.MType, metric.ID, strconv.FormatFloat(*metric.Value, 'f', -1, 64))
	_, err := http.Post(url, "application/json", nil)
	if err != nil {
		logger.Error("Failed to write metric:")
		logger.Error(err.Error())
	} else {
		logger.Log(fmt.Sprintf("Successfully updated metric %s", metric.ID))
	}
}
