package api

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
)

// Server is a minimal HTTP server exposing health and metrics endpoints.
type Server struct {
	addr   string
	ready  atomic.Bool
	server *http.Server
}

// New creates a new Server listening on addr (e.g. ":8080").
func New(addr string) *Server {
	s := &Server{addr: addr}

	// Set HTTP Router
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleLiveness)
	mux.HandleFunc("/readyz", s.handleReadiness)
	mux.HandleFunc("/metrics", s.handleMetrics)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s
}

// SetReady marks the agent as ready (true) or not ready (false).
func (s *Server) SetReady(v bool) {
	s.ready.Store(v)
}

// Start begins listening in a background goroutine.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	go func() { _ = s.server.Serve(ln) }()
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) handleLiveness(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleReadiness(w http.ResponseWriter, _ *http.Request) {
	if s.ready.Load() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("not ready"))
	}
}

func (s *Server) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	// Placeholder: will be replaced by Prometheus exposition format.
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("# perf-agent metrics placeholder\n"))
}
