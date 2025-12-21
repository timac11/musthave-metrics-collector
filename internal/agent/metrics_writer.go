package agent

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type MetricsWriterConfig struct {
	Attempts         uint
	AttemptsInterval uint
}

type MetricsWriter struct {
	client resty.Client
	config MetricsWriterConfig
}

func (mw *MetricsWriter) Write(metrics []model.Metrics) {
	var res *resty.Response
	var err error

	err = retry.Do(
		func() error {
			res, err = mw.client.R().SetBody(metrics).Post("/updates")
			return err
		},
		mw.getRetryOptions()...,
	)

	if err != nil {
		logger.Error("Failed to write metrics", err.Error())
		return
	}

	if res.StatusCode() != http.StatusOK {
		logger.Error("Failed to write metrics", "status", res.StatusCode())
		return
	}

	logger.Info("Success updated metrics", "status", res.Status())
}

func (mw *MetricsWriter) writeMetric(metric model.Metrics) {
	var res *resty.Response
	var err error

	err = retry.Do(
		func() error {
			res, err = mw.client.R().SetBody(metric).Post("/update")
			return err
		},
		mw.getRetryOptions()...,
	)

	if err != nil {
		logger.Error("Failed to write metric", "id", metric.ID, "type", metric.MType)
		logger.Error(err.Error())
		return
	}

	logger.Info("Success update metric", "metric ID", metric.ID, "status", res.Status())
}

func newMetricsWriter(url string, config MetricsWriterConfig) *MetricsWriter {
	client := resty.New()

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	client.SetBaseURL(url)

	mw := &MetricsWriter{
		client: *client,
		config: config,
	}

	return mw
}

func (mw *MetricsWriter) getRetryOptions() []retry.Option {
	return []retry.Option{
		retry.Attempts(mw.config.Attempts),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return time.Second + time.Duration(n*mw.config.AttemptsInterval)*time.Second
		}),
		retry.Context(context.Background()),
	}
}
