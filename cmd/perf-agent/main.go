package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paoloelefante/perf-agent/internal/api"
	"github.com/paoloelefante/perf-agent/internal/version"
)

func main() {
	slog.Info("perf-agent starting", "version", version.Version)

	healthAddr := os.Getenv("HEALTH_ADDR")
	if healthAddr == "" {
		healthAddr = ":8080"
	}

	srv := api.New(healthAddr)
	if err := srv.Start(); err != nil {
		slog.Error("failed to start api server", "err", err)
		os.Exit(1)
	}

	srv.SetReady(true)
	slog.Info("perf-agent ready", "addr", healthAddr)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	slog.Info("perf-agent shutting down")

	srv.SetReady(false)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("api server shutdown error", "err", err)
	}

	slog.Info("perf-agent stopped")
}
