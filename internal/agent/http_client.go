package agent

import (
	"strings"
	"time"

	"github.com/timac11/musthave-metrics-collector/internal/common/encryption"
	"github.com/timac11/musthave-metrics-collector/internal/common/util"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"

	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type HTTPClient struct {
	client     resty.Client
	encoder    *encryption.Encoder
	signingKey string
}

type HTTPClientConfig struct {
	address    string
	cryptoKey  string
	signingKey string
}

func newHTTPClient(conf HTTPClientConfig) (*HTTPClient, error) {
	client := resty.New()
	url := conf.address

	if !strings.HasPrefix(conf.address, "http") {
		url = "http://" + url
	}

	client.SetBaseURL(url)
	client.SetTimeout(time.Duration(10 * time.Second))

	var encoder *encryption.Encoder

	if conf.cryptoKey != "" {
		newEncoder, err := encryption.NewEncoder(conf.cryptoKey)

		if err != nil {
			return nil, err
		}

		encoder = newEncoder
	}

	return &HTTPClient{client: *client, signingKey: conf.signingKey, encoder: encoder}, nil
}

func (client *HTTPClient) calculateRequestBody(body any) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return []byte{}, err
	}

	if client.encoder != nil {
		return client.encoder.Encode(data)
	}

	return data, nil
}

func (client *HTTPClient) UpdateMetrics(ctx context.Context, metrics []model.Metrics) error {
	var err error

	body, err := client.calculateRequestBody(metrics)
	if err != nil {
		return err
	}

	signature, err := util.CalculateSignature(body, client.signingKey)

	if err != nil {
		logger.Error("Failed to calculate signature", err.Error())
		return err
	}

	_, err = client.client.R().SetBody(body).SetHeader("HashSHA256", signature).Post("/updates")
	return err
}

func (client *HTTPClient) UpdateMetric(ctx context.Context, metric model.Metrics) error {
	var err error

	body, err := client.calculateRequestBody(metric)
	if err != nil {
		return err
	}

	signature, err := util.CalculateSignature(body, client.signingKey)

	if err != nil {
		logger.Error("Failed to calculate signature", err.Error())
		return err
	}

	_, err = client.client.R().SetBody(body).SetHeader("HashSHA256", signature).Post("/update")
	return err
}
