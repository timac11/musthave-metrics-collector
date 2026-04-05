package grpc

import (
	"context"

	"github.com/timac11/musthave-metrics-collector/internal/model"
	"github.com/timac11/musthave-metrics-collector/internal/service"

	pb "github.com/timac11/musthave-metrics-collector/internal/grpc/proto"
)

type GrpcServer struct {
	pb.UnimplementedMetricsServer
	service service.Service
}

func (server *GrpcServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var response pb.UpdateMetricsResponse
	var modelMetrics []*model.Metrics

	for _, metric := range in.GetMetrics() {
		mtype := metric.GetType()
		modelMetric := model.Metrics{
			ID: metric.GetId(),
		}

		if mtype == 0 {
			value := metric.GetValue()
			modelMetric.MType = model.Gauge
			modelMetric.Value = &value
		} else {
			delta := metric.GetDelta()
			modelMetric.MType = model.Counter
			modelMetric.Delta = &delta
		}

		modelMetrics = append(modelMetrics, &modelMetric)
	}

	server.service.SaveAll(ctx, modelMetrics)

	return &response, nil
}
