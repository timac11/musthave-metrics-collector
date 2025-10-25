package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"

	logger "github.com/timac11/musthave-metrics-collector/internal/logger"
	model "github.com/timac11/musthave-metrics-collector/internal/model"
)

type MetricsWriter struct {
	client resty.Client
}

func (mw *MetricsWriter) Write(metrics []model.Metrics) {
	for _, metric := range metrics {
		mw.writeMetric(metric)
	}
}

func (mw *MetricsWriter) writeMetric(metric model.Metrics) {
	value := strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	url := fmt.Sprintf("update/%s/%s/%s", metric.MType, metric.ID, value)
	res, err := mw.client.R().Post(url)
	if err != nil {
		logger.Error("Failed to write metric:")
		logger.Error(err.Error())
		return
	}

	logger.Log(fmt.Sprintf("updated metric %s status %s", metric.ID, res.Status()))
}

func NewMetricsWriter(url string) *MetricsWriter {
	client := resty.New()
	client.SetBaseURL(url)
	client.SetScheme("http")

	mw := &MetricsWriter{
		client: *client,
	}

	return mw
}
