package handler

import (
	"context"
	"log/slog"

	"github.com/alexliesenfeld/health"

	"github.com/mapout-world/stern/lib/bedrock"

	placespb "github.com/mapout-world/stern/services/service.places/proto/v1"
)

var _ placespb.PlacesServer = (*Handler)(nil)
var _ bedrock.ServiceHandler = (*Handler)(nil)

type Handler struct {
	placespb.UnimplementedPlacesServer

	log     *slog.Logger
	metrics Metrics
}

func New(log *slog.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) HealthChecks() []health.Check {
	return []health.Check{}
}

func (h *Handler) GetPlace(ctx context.Context, req *placespb.GetPlaceRequest) (*placespb.GetPlaceReply, error) {
	h.metrics.ExampleCounter.Inc()
	return &placespb.GetPlaceReply{}, nil
}
