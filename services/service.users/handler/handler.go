package handler

import (
	"context"
	"log/slog"

	"github.com/alexliesenfeld/health"

	"github.com/mapout-world/stern/lib/bedrock"

	userspb "github.com/mapout-world/stern/services/service.users/proto/v1"
)

var _ userspb.UsersServer = (*Handler)(nil)
var _ bedrock.ServiceHandler = (*Handler)(nil)

type Handler struct {
	userspb.UnimplementedUsersServer

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

func (h *Handler) GetUser(ctx context.Context, req *userspb.GetUserRequest) (*userspb.GetUserReply, error) {
	h.metrics.ExampleCounter.Inc()
	return &userspb.GetUserReply{}, nil
}
