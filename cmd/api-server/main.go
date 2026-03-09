package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ihyaulhaq/url-shotener-BE/internal/config"
	"github.com/ihyaulhaq/url-shotener-BE/internal/handler"
	"github.com/ihyaulhaq/url-shotener-BE/internal/middleware"
)

func main() {
	const PORT = "8080"

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
	}

	h := handler.New()

	chain := middleware.Chaining()

	server := &http.Server{
		Addr:    PORT,
		Handler: chain(h.Routes()),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", server.Addr, "env", cfg.App.Env)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")

}
