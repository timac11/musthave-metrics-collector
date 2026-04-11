package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/timac11/musthave-metrics-collector/internal/common/util"
	pb "github.com/timac11/musthave-metrics-collector/internal/grpc/proto"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type GrpcClient struct {
	client  pb.MetricsClient
	localIp string
}

type GrpcClientConfig struct {
	Address string
}

// NewGrpcClient constructor
func NewGrpcClient(conf GrpcClientConfig) (*GrpcClient, error) {
	localIp, err := util.GetOutboundIP(conf.Address)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(
		conf.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(clientIPInterceptor(localIp)),
	)

	if err != nil {
		return nil, err
	}

	client := pb.NewMetricsClient(conn)

	return &GrpcClient{client: client, localIp: localIp}, nil
}

// UpdateMetrics is wrapper on grpc client method of batch update
func (client *GrpcClient) UpdateMetrics(ctx context.Context, metrics []model.Metrics) error {
	var sendMetrics []*pb.Metric

	for _, item := range metrics {
		metric := pb.Metric_builder{Id: item.ID, Value: *item.Value, Delta: *item.Delta, Type: pb.Metric_MType(*item.Delta)}.Build()
		sendMetrics = append(sendMetrics, metric)
	}

	in := pb.UpdateMetricsRequest_builder{
		Metrics: sendMetrics,
	}.Build()

	md := metadata.Pairs(
		"x-real-ip", client.localIp,
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := client.client.UpdateMetrics(ctx, in)

	return err
}

// UpdateMetrics is wrapper on grpc client method of single update
func (client *GrpcClient) UpdateMetric(ctx context.Context, metric model.Metrics) error {
	sendMetrics := []model.Metrics{
		metric,
	}

	return client.UpdateMetrics(ctx, sendMetrics)
}
