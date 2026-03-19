package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/not-kamalesh/pismo-account/api"
	"github.com/not-kamalesh/pismo-account/internal/healthcheck"
)

func main() {
	// Init Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load Config @TODO

	// Init Database @TODO

	// Init Clients @TODO
	healthCheckHandler := healthcheck.NewHandler()

	apiHandler := api.NewAPIHandler(healthCheckHandler)

	// Register Routes and Handlers
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("").Subrouter()
	apiRouter.HandleFunc("/health_check", apiHandler.HealthCheck).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrChan := make(chan error, 1)

	// Start server
	go func() {
		slog.Info("http server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()

	// Handle inturrupts and termination
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		slog.Warn("received shutdown signal")
	case err := <-serverErrChan:
		slog.Error("server failed unexpectedly", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	} else {
		slog.Info("server shutdown gracefully")
	}
}
