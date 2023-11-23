package bedrock

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/alexliesenfeld/health"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/mapout-world/stern/lib/healthcheck"
	"github.com/mapout-world/stern/lib/metrics"
)

var _ InstrumentedServer = (*GRPCServer)(nil)

type GRPCServer struct {
	srv          *grpc.Server
	healthchecks []health.Check

	addr string
}

func NewGRPCServer(log *slog.Logger, port int) *GRPCServer {
	grpcLogger := func(l *slog.Logger) logging.Logger {
		return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		})
	}(log)

	grpcMetrics := metrics.NewGRPCMetrics()

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(
				logging.UnaryServerInterceptor(grpcLogger,
					logging.WithLevels(func(code codes.Code) logging.Level {
						return logging.LevelDebug
					}),
				),
			),
			grpc.UnaryServerInterceptor(
				grpcMetrics.UnaryServerInterceptor(
					grpcprom.WithExemplarFromContext(metrics.ExemplarFromContext),
				),
			),
		),
	}

	addr := fmt.Sprintf(":%d", port)
	return &GRPCServer{
		srv: grpc.NewServer(opts...),
		healthchecks: []health.Check{
			healthcheck.GRPCCheck("bedrock", addr),
		},
		addr: addr,
	}
}

type ServiceHandler interface {
	HealthChecks() []health.Check
	RegisterMetrics()
}

func (s *GRPCServer) RegisterService(sd *grpc.ServiceDesc, h ServiceHandler) {
	s.srv.RegisterService(sd, h)
	s.healthchecks = append(s.healthchecks, h.HealthChecks()...)
	h.RegisterMetrics()
}

func (s *GRPCServer) HealthChecks() []health.Check {
	return s.healthchecks
}

func (s *GRPCServer) Name() string {
	return "grpc"
}

func (s *GRPCServer) Addr() string {
	return s.addr
}

func (s *GRPCServer) Serve(ctx context.Context) error {
	reflection.Register(s.srv)
	healthpb.RegisterHealthServer(s.srv, &HealthHandler{})

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc server failed to serve: %w", err)
	}

	return s.srv.Serve(lis)
}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	ok := make(chan struct{})

	go func() {
		s.srv.GracefulStop()
		close(ok)
	}()

	select {
	case <-ok:
		return nil
	case <-ctx.Done():
		s.srv.Stop()
		return ctx.Err()
	}
}

type HealthHandler struct {
}

func (h *HealthHandler) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthHandler) Watch(req *healthpb.HealthCheckRequest, server healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
