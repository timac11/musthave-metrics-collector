package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"

	"github.com/timac11/musthave-metrics-collector/internal/common/util"
	"github.com/timac11/musthave-metrics-collector/internal/encryption"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

// MetricsWriterConfig is the configuration of metric writer
type MetricsWriterConfig struct {
	Attempts         uint   // number of attempts to resend metrics
	AttemptsInterval uint   // parameter for calculation of backoff interval between two attempts. backoff time on i-th iteration is equal ti (i-1) * AttemptsInterval
	SigningKey       string // used to sign metrics before sending it to the server
	CryptoKey        string // used for encryption metrics payload
}

// MetricsWriter is structure of writer
type MetricsWriter struct {
	client  resty.Client // client is used to send metrics to server
	config  MetricsWriterConfig
	encoder *encryption.Encoder // encoder is used for metrics encryption
}

// Write used to send metrics to server
func (mw *MetricsWriter) Write(metrics []model.Metrics) error {
	var res *resty.Response
	var err error

	body, err := mw.calculateRequestBody(metrics)
	if err != nil {
		return err
	}

	signature, err := util.CalculateSignature(body, mw.config.SigningKey)

	if err != nil {
		return err
	}

	err = retry.Do(
		func() error {
			res, err = mw.client.R().SetBody(body).SetHeader("HashSHA256", signature).Post("/updates")
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

func (mw *MetricsWriter) calculateRequestBody(body any) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return []byte{}, err
	}

	if mw.encoder != nil {
		return mw.encoder.Encode(data)
	}

	return data, nil
}

func (mw *MetricsWriter) writeMetric(metric model.Metrics) error {
	var res *resty.Response
	var err error

	body, err := mw.calculateRequestBody(metric)
	if err != nil {
		return err
	}

	signature, err := util.CalculateSignature(body, mw.config.SigningKey)

	if err != nil {
		logger.Error("Failed to calculate signature", err.Error())
		return err
	}

	err = retry.Do(
		func() error {
			res, err = mw.client.R().SetBody(body).SetHeader("HashSHA256", signature).Post("/update")
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

func newMetricsWriter(url string, config MetricsWriterConfig) (*MetricsWriter, error) {
	client := resty.New()

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	client.SetBaseURL(url)

	if config.CryptoKey != "" {
		encoder, err := encryption.NewEncoder(config.CryptoKey)

		if err != nil {
			return nil, err
		}

		return &MetricsWriter{
			client:  *client,
			config:  config,
			encoder: encoder,
		}, nil
	}

	return &MetricsWriter{
		client: *client,
		config: config,
	}, nil
}
