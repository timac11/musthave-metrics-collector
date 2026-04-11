package grpc

import (
	"context"
	"net/netip"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ServerInterceptor struct {
	subnet *netip.Prefix
}

func (s *ServerInterceptor) ServerIPInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if s.subnet != nil {
		var realIP string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("x-real-ip")
			if len(values) > 0 {
				realIP = values[0]
			}
		}

		ip, err := netip.ParseAddr(realIP)

		if err != nil {
			logger.Error("Failed to check client ip", realIP)
			return nil, status.Error(codes.Unauthenticated, "Client ip is not in trusted subnet")
		}

		if !s.subnet.Contains(ip) {
			logger.Error("Client ip is not in trusted subnet", realIP)
			return nil, status.Error(codes.Unauthenticated, "Client ip is not in trusted subnet")
		}

	}

	return handler(ctx, req)
}

func (s *ServerInterceptor) LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	logger.Info("GRPC Request info", info.FullMethod)
    resp, err := handler(ctx, req)

	if err != nil {
		logger.Error("GRPC Request error", info.FullMethod, err)
	} else {
		logger.Info("GRPC Request succeeded", info.FullMethod)
	}

    return resp, err
}


func NewServerInterceptor(subnet string) (*ServerInterceptor, error) {
	if subnet != "" {
		parsed, err := netip.ParsePrefix(subnet)

		if err != nil {
			return nil, err
		}

		return &ServerInterceptor{subnet: &parsed}, nil
	}

	return &ServerInterceptor{}, nil
}
