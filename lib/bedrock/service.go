package bedrock

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	log     *slog.Logger
	servers []Server
}

func NewService(log *slog.Logger, servers ...Server) *Service {
	return &Service{
		log:     log,
		servers: servers,
	}
}

type Server interface {
	Name() string
	Addr() string
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

func (s *Service) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, srv := range s.servers {
		srv := srv
		g.Go(func() error {
			s.log.Info("server listening", "name", srv.Name(), "addr", srv.Addr())
			return srv.Serve(ctx)
		})
	}

	s.log.Info("service started")
	<-ctx.Done()
	s.log.Info("service shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, srv := range s.servers {
		srv := srv
		g.Go(func() error {
			if err := srv.Shutdown(ctx); err != nil {
				s.log.Warn("server shutdown error", "name", srv.Name(), "error", err)
			} else {
				s.log.Info("server shutdown complete", "name", srv.Name())
			}
			return nil
		})
	}

	defer s.log.Info("service stopped")
	return g.Wait()
}
