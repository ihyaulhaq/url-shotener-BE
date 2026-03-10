package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ihyaulhaq/url-shotener-BE/internal/config"
	"github.com/ihyaulhaq/url-shotener-BE/internal/handler"
	"github.com/ihyaulhaq/url-shotener-BE/internal/middleware"
	"github.com/ihyaulhaq/url-shotener-BE/internal/service"
	"github.com/ihyaulhaq/url-shotener-BE/internal/store"

	_ "github.com/lib/pq"
)

func main() {

	// Load all config server, db, app
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	//open DB
	db, err := openDB(cfg.DB)
	if err != nil {
		slog.Error("Failed to connect to db", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Db connected:", "host", cfg.DB.Host, "name", cfg.DB.Name)

	s := store.NewStore(db)
	urlSrv := service.NewUrlService(s)
	h := handler.New(urlSrv, cfg.App.BaseURL)

	chain := middleware.Chaining()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      chain(h.Routes()),
		ReadTimeout:  cfg.Server.ReadTimeOut,
		WriteTimeout: cfg.Server.WriteTimeOut,
		IdleTimeout:  cfg.Server.IdleTime,
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

func openDB(cfg config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open db :%w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}
