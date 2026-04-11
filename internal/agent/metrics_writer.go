package agent

import (
	"context"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/timac11/musthave-metrics-collector/internal/agent/grpc"
	"github.com/timac11/musthave-metrics-collector/internal/agent/http"
	"github.com/timac11/musthave-metrics-collector/internal/common/encryption"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

// MetricsWriterConfig is the configuration of metric writer
type MetricsWriterConfig struct {
	URL              string // server address
	Mode             string // http or grpc transport mode
	Attempts         uint   // number of attempts to resend metrics
	AttemptsInterval uint   // parameter for calculation of backoff interval between two attempts. backoff time on i-th iteration is equal ti (i-1) * AttemptsInterval
	SigningKey       string // used to sign metrics before sending it to the server
	CryptoKey        string // used for encryption metrics payload
}

// Client is interface to send metrics to server
// in this case it is grpc or http client
type Client interface {
	UpdateMetrics(ctx context.Context, metrics []model.Metrics) error
	UpdateMetric(ctx context.Context, metric model.Metrics) error
}

// MetricsWriter is structure of writer
type MetricsWriter struct {
	client  Client // client is used to send metrics to server
	config  MetricsWriterConfig
	encoder *encryption.Encoder // encoder is used for metrics encryption
}

// Write used to send metrics to server
func (mw *MetricsWriter) Write(ctx context.Context, metrics []model.Metrics) error {
	err := retry.Do(
		func() error {
			err := mw.client.UpdateMetrics(ctx, metrics)
			return err
		},
		// mw.getRetryOptions()...,
	)

	if err != nil {
		return err
	}

	logger.Info("Success updated metrics")
	return nil
}

func (mw *MetricsWriter) writeMetric(ctx context.Context, metric model.Metrics) error {
	err := retry.Do(
		func() error {
			err := mw.client.UpdateMetric(ctx, metric)
			if err != nil {
				logger.Error(err.Error())
			}
			return err
		},
		mw.getRetryOptions()...,
	)

	if err != nil {
		logger.Error("Failed to write metric", "id", metric.ID, "type", metric.MType, err.Error())
		return err
	}

	logger.Info("Success update metric", "metric ID", metric.ID)
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

func newMetricsWriter(config MetricsWriterConfig) (*MetricsWriter, error) {
	var client Client

	if config.Mode == "http" {
		client, err := http.NewHTTPClient(http.HTTPClientConfig{Address: config.URL, CryptoKey: config.CryptoKey, SigningKey: config.SigningKey})

		if err != nil {
			return nil, err
		}

		return &MetricsWriter{client: client, config: config}, nil
	}

	client, err := grpc.NewGrpcClient(grpc.GrpcClientConfig{Address: config.URL})

	if err != nil {
		return nil, err
	}

	return &MetricsWriter{client: client, config: config}, nil
}
