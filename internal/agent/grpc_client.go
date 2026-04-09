package agent

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/timac11/musthave-metrics-collector/internal/grpc/proto"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type GrpcClient struct {
	client pb.MetricsClient
}

type GrpcClientConfig struct {
	address string
}

func newGrpcClient(conf GrpcClientConfig) (*GrpcClient, error) {
	conn, err := grpc.NewClient(conf.address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := pb.NewMetricsClient(conn)

	return &GrpcClient{client: client}, nil
}

func (client *GrpcClient) UpdateMetrics(ctx context.Context, metrics []model.Metrics) error {
	var sendMetrics []*pb.Metric

	for _, item := range metrics {
		metric := pb.Metric_builder{Id: item.ID, Value: *item.Value, Delta: *item.Delta, Type: pb.Metric_MType(*item.Delta)}.Build()
		sendMetrics = append(sendMetrics, metric)
	}

	in := pb.UpdateMetricsRequest_builder{
		Metrics: sendMetrics,
	}.Build()

	_, err := client.client.UpdateMetrics(ctx, in)

	return err
}

func (client *GrpcClient) UpdateMetric(ctx context.Context, metric model.Metrics) error {
	sendMetrics := []model.Metrics{
		metric,
	}

	return client.UpdateMetrics(ctx, sendMetrics)
}
