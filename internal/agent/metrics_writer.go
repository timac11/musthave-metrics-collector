package agent

import (
	"github.com/go-resty/resty/v2"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"strings"
)

type MetricsWriter struct {
	client resty.Client
}

func (mw *MetricsWriter) Write(metrics []model.Metrics) {
	res, err := mw.client.R().SetBody(metrics).Post("/updates")

	if err != nil {
		logger.Error("Failed to write metrics")
		logger.Error(err.Error())
		return
	}

	if res != nil && res.StatusCode() != 200 {
		logger.Error("Failed to write metrics", "status", res.StatusCode())
		return
	}

	logger.Info("Success updated metrics", "status", res.Status())
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

	client.SetRetryCount(3)

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	client.SetBaseURL(url)

	mw := &MetricsWriter{
		client: *client,
	}

	return mw
}
