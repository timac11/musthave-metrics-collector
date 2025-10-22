package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"sync"

	logger "github.com/timac11/musthave-metrics-collector/internal/logger"
	model "github.com/timac11/musthave-metrics-collector/internal/model"
)

var baseURL = "http://localhost:8080"

var (
	client *resty.Client
	once   sync.Once
)

func Write(metrics []model.Metrics) {
	for _, metric := range metrics {
		writeMetric(metric)
	}
}

func getClient() *resty.Client {
	once.Do(func() {
		client = resty.New()
		client.SetBaseURL(baseURL)
	})
	return client
}

func writeMetric(metric model.Metrics) {
	value := strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	url := fmt.Sprintf("%s/update/%s/%s/%s", baseURL, metric.MType, metric.ID, value)
	res, err := getClient().R().Post(url)
	if err != nil {
		logger.Error("Failed to write metric:")
		logger.Error(err.Error())
		return
	}

	logger.Log(fmt.Sprintf("updated metric %s status %s", metric.ID, res.Status()))
}
