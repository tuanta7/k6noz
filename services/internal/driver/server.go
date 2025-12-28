package driver

import (
	"context"
	"net/http"
)

type Server struct {
	mux     *http.ServeMux
	server  *http.Server
	handler *Handler
}

func NewServer(addr string, handler *Handler) *Server {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &Server{
		mux:     mux,
		server:  server,
		handler: handler,
	}
}

func (s *Server) Run() error {
	s.mux.HandleFunc("GET /drivers/{id}", s.handler.GetDriverByID)
	s.mux.HandleFunc("POST /ratings", s.handler.CreateNewRating)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
