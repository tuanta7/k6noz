package ingestion

import (
	"context"
	"net/http"

	"github.com/tuanta7/k6noz/services/pkg/otelx"
)

type Server struct {
	server     *http.Server
	mux        *http.ServeMux
	handler    *Handler
	prometheus *otelx.PrometheusProvider
}

func NewServer(cfg *Config, handler *Handler, prometheus *otelx.PrometheusProvider) *Server {
	mux := http.NewServeMux()

	return &Server{
		handler:    handler,
		prometheus: prometheus,
		mux:        mux,
		server: &http.Server{
			Addr:    cfg.BindAddress,
			Handler: mux,
		},
	}
}

func (s *Server) Run() error {
	s.registerRoutes()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	s.mux.Handle("GET /metrics", s.prometheus.Handler())
	s.mux.HandleFunc("/ws", s.handler.HandleWS)
}
