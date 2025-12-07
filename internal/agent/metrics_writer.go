package agent

import (
	"github.com/go-resty/resty/v2"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"strings"
	"time"
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
	res, err := mw.client.R().SetBody(metric).Post("/update")

	if err != nil {
		logger.Error("Failed to write metric", "id", metric.ID, "type", metric.MType)
		logger.Error(err.Error())
		return
	}

	logger.Info("Success update metric", "metric ID", metric.ID, "status", res.Status())
}

func newMetricsWriter(url string) *MetricsWriter {
	client := resty.New()

	client.SetRetryWaitTime(1* time.Second).SetRetryCount(10)

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	client.SetBaseURL(url)

	mw := &MetricsWriter{
		client: *client,
	}

	return mw
}
