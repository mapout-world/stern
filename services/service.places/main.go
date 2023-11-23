package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mapout-world/stern/lib/bedrock"
	"github.com/mapout-world/stern/lib/config"
	"github.com/mapout-world/stern/lib/logging"
	"github.com/mapout-world/stern/services/service.places/handler"

	placespb "github.com/mapout-world/stern/services/service.places/proto/v1"
)

var cfg struct {
	config.StaticBase
}

func main() {
	config.MustLoadStatic(&cfg)
	config.MustLoadDynamic(cfg.Config)
	defer config.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop()
	}()

	logger := logging.NewLogger(ctx, os.Stdout)
	h := handler.New(logger)

	srv := bedrock.NewGRPCServer(logger, cfg.Port)
	srv.RegisterService(&placespb.Places_ServiceDesc, h)

	admin := bedrock.NewAdminServer(cfg.AdminPort)
	admin.Instrument(srv)

	svc := bedrock.NewService(logger, srv, admin)
	if err := svc.Run(ctx); err != nil {
		logger.Error("error running service", "error", err)
		os.Exit(1)
	}
}
