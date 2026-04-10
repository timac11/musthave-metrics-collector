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

func NewServer(addr string, service service.Service) *Server {
	server := grpc.NewServer()
	return &Server{grpcServer: server, addr: addr, service: service}
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
