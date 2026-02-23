package agent

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"

	"github.com/timac11/musthave-metrics-collector/internal/common/util"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

// MetricsWriterConfig is the configuration of metric writer
type MetricsWriterConfig struct {
	Attempts         uint   // number of attempts to resend metrics
	AttemptsInterval uint   // parameter for calculation of backoff interval between two attempts. backoff time on i-th iteration is equal ti (i-1) * AttemptsInterval
	SigningKey       string // used to sign metrics before sending it to the server
}

// MetricsWriter is structure of writer
type MetricsWriter struct {
	client resty.Client // client is used to send metrics to server
	config MetricsWriterConfig
}

// Write used to send metrics to server
func (mw *MetricsWriter) Write(metrics []model.Metrics) error {
	var res *resty.Response
	var err error

	signature, err := util.CalculateSignature(metrics, mw.config.SigningKey)

	if err != nil {
		return err
	}

	err = retry.Do(
		func() error {
			res, err = mw.client.R().SetBody(metrics).SetHeader("HashSHA256", signature).Post("/updates")
			return err
		},
		mw.getRetryOptions()...,
	)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to write metrics, status = %d", res.StatusCode())
	}

	logger.Info("Success updated metrics", "status", res.Status())
	return nil
}

func (mw *MetricsWriter) writeMetric(metric model.Metrics) error {
	var res *resty.Response
	var err error
	signature, err := util.CalculateSignature(metric, mw.config.SigningKey)

	if err != nil {
		logger.Error("Failed to calculate signature", err.Error())
		return err
	}

	err = retry.Do(
		func() error {
			res, err = mw.client.R().SetBody(metric).SetHeader("HashSHA256", signature).Post("/update")
			return err
		},
		mw.getRetryOptions()...,
	)

	if err != nil {
		logger.Error("Failed to write metric", "id", metric.ID, "type", metric.MType, err.Error())
		return err
	}

	logger.Info("Success update metric", "metric ID", metric.ID, "status", res.Status())
	return nil
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
