package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "github.com/timac11/musthave-metrics-collector/internal/grpc/proto"
	"github.com/timac11/musthave-metrics-collector/internal/grpc/server"
	"github.com/timac11/musthave-metrics-collector/internal/service"
)

type Server struct {
	grpcServer *grpc.Server
	addr       string
	service    service.Service
}

func NewServer(addr, subnet string, service service.Service) (*Server, error) {
	interceptor, err := NewServerInterceptor(subnet)

	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor.ServerIPInterceptor))
	return &Server{grpcServer: server, addr: addr, service: service}, nil
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	pb.RegisterMetricsServer(s.grpcServer, server.NewMetricServer(s.service))

	return s.grpcServer.Serve(listen)
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}
