package agent

import (
	"fmt"
	"net/http"
	"strconv"

	logger "github.com/timac11/musthave-metrics-collector/internal/logger"
	model "github.com/timac11/musthave-metrics-collector/internal/model"
)

var baseURL = "http://localhost:8080"

func Write(metrics []model.Metrics) {
	for _, metric := range metrics {
		writeMetric(metric)
	}
}

func writeMetric(metric model.Metrics) {
	url := fmt.Sprintf("%s/%s/%s/%s", baseURL, metric.MType, metric.ID, strconv.FormatFloat(*metric.Value, 'f', -1, 64))
	res, err := http.Post(url, "application/json", nil)
	if err != nil {
		logger.Error("Failed to write metric:")
		logger.Error(err.Error())
		return
	}

	logger.Log(fmt.Sprintf("updated metric %s status %s", metric.ID, res.Status))
	defer res.Body.Close()
}
