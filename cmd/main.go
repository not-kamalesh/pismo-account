package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/not-kamalesh/pismo-account/api"
	"github.com/not-kamalesh/pismo-account/internal/account"
	"github.com/not-kamalesh/pismo-account/internal/healthcheck"
	"github.com/not-kamalesh/pismo-account/internal/idempotencymgr"
	"github.com/not-kamalesh/pismo-account/internal/transaction"
	"github.com/not-kamalesh/pismo-account/server"
	"github.com/not-kamalesh/pismo-account/storage"
)

func main() {
	// Init Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Load Config from the file
	appConfig, err := server.LoadConfig()
	if err != nil {
		// if load config fails, app can't run properly, so panic
		panic(fmt.Sprintf("error occured on loading app config, err= %v", err))
	}

	// Init the client objects required for the server (ex: DB, Redis ...)
	clients, err := server.InitClients(appConfig)
	if err != nil {
		// if the primary component initialisation fails, app can't run properly, so panic
		panic(fmt.Sprintf("error occured on InitClients, err= %v", err))
	}

	// Init the DAOs
	accountDao := storage.NewAccountDao(clients.DB)
	transactionDao := storage.NewTransactionDao(clients.DB)

	// Init Handlers
	healthCheckHandler := healthcheck.NewHandler()
	accountHandler := account.NewHandler(accountDao)
	transctionHandler := transaction.NewHandler(accountDao, transactionDao)
	idempotencyMgr := idempotencymgr.NewInMemIdempotencyMgr()
	apiHandler := api.NewAPIHandler(healthCheckHandler, accountHandler, transctionHandler, idempotencyMgr)

	// Register Routes and Handlers
	r := mux.NewRouter()

	apiRouter := r.PathPrefix("").Subrouter()
	apiRouter.Use(recoveryMiddleware) // api panic recovary middleware
	apiRouter.HandleFunc("/health_check", apiHandler.HealthCheck).Methods(http.MethodGet)
	apiRouter.HandleFunc("/accounts", apiHandler.CreateAccount).Methods(http.MethodPost)
	apiRouter.HandleFunc("/accounts/{account_id}", apiHandler.GetAccount).Methods(http.MethodGet)
	apiRouter.HandleFunc("/transactions", apiHandler.CreateTransaction).Methods(http.MethodPost)

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

// middleware to recovers from panics
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic occurred", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
