package bedrock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alexliesenfeld/health"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mapout-world/stern/lib/config"
	"github.com/mapout-world/stern/lib/healthcheck"
)

type AdminServer struct {
	srv    *http.Server
	health health.Checker
}

func NewAdminServer(port int) *AdminServer {
	return &AdminServer{
		srv: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		health: healthcheck.NewChecker(),
	}
}

func (s *AdminServer) Name() string {
	return "admin"
}

func (s *AdminServer) Addr() string {
	return s.srv.Addr
}

type InstrumentedServer interface {
	HealthChecks() []health.Check
}

func (s *AdminServer) Instrument(target InstrumentedServer) {
	s.health = healthcheck.NewChecker(target.HealthChecks()...)
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
}

func (s *AdminServer) Serve(ctx context.Context) error {
	s.srv.Handler = s.router()

	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("admin server failed to serve: %w", err)
	}

	return nil
}

func (s *AdminServer) router() http.Handler {
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	router.HandleFunc("/ping", s.HandlePing).Methods(http.MethodGet)
	router.HandleFunc("/config", s.HandleConfig).Methods(http.MethodGet)

	subrouter := router.PathPrefix("/health").Subrouter()
	subrouter.Handle("/startup", health.NewHandler(s.health)).Methods(http.MethodGet)
	subrouter.Handle("/live", health.NewHandler(s.health)).Methods(http.MethodGet)
	subrouter.Handle("/ready", health.NewHandler(s.health)).Methods(http.MethodGet)

	return router
}

func (s *AdminServer) HandlePing(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(map[string]string{
		config.Get("name").String("unknown"): "OK",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *AdminServer) HandleConfig(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(config.All())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *AdminServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
