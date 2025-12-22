package agent

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	SigningKey       string
}

type MetricsWriter struct {
	client resty.Client
	config MetricsWriterConfig
}

func (mw *MetricsWriter) Write(metrics []model.Metrics) error {
	var res *resty.Response
	var err error

	signature, err := mw.getSignuture(metrics)

	if err != nil {
		logger.Error("Failed to calculate signature", err.Error())
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
		logger.Error("Failed to write metrics", err.Error())
		return err
	}

	if res.StatusCode() != http.StatusOK {
		logger.Error("Failed to write metrics", "status", res.StatusCode())
		return err
	}

	logger.Info("Success updated metrics", "status", res.Status())
	return nil
}

func (mw *MetricsWriter) writeMetric(metric model.Metrics) error {
	var res *resty.Response
	var err error
	signature, err := mw.getSignuture(metric)

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

func (mw *MetricsWriter) getSignuture(obj interface{}) (string, error) {
	if mw.config.SigningKey == "" {
		return "", nil
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	hashInBytes := sha256.Sum256(bytes)
	return hex.EncodeToString(hashInBytes[:]), nil
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
