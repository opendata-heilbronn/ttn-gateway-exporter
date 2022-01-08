package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	server := &Server{
		server: httpServer,
	}
	return server
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}
