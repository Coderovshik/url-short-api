package main

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/Coderovshik/url-short/internal/config"
	"github.com/Coderovshik/url-short/internal/http-server/handlers/ping"
	"github.com/Coderovshik/url-short/internal/http-server/handlers/url/delete"
	"github.com/Coderovshik/url-short/internal/http-server/handlers/url/redirect"
	"github.com/Coderovshik/url-short/internal/http-server/handlers/url/save"
	"github.com/Coderovshik/url-short/internal/http-server/middleware/loggermw"
	"github.com/Coderovshik/url-short/internal/lib/logger/slut"
	"github.com/Coderovshik/url-short/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	logText = "text"
	logJSON = "json"

	levelDebug = "debug"
	levelInfo  = "info"

	outStd = "std"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs/config.json", "config path")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := setupLogger(&cfg.Logger)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	connectionString := os.Getenv("CONNECTION_STRING")
	if len(connectionString) == 0 {
		log.Error("failed to retrieve env variable: CONNECTION_STRING")
		os.Exit(1)
	}

	myCtx := context.Background()
	storage, err := postgres.New(myCtx, connectionString)
	if err != nil {
		log.Error("failed to init storage", slut.ErrAttr(err))
		os.Exit(1)
	}
	defer storage.Close()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(loggermw.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/ping", ping.New(log))
	router.Post("/url", save.New(myCtx, log, storage))
	router.Get("/{alias}", redirect.New(myCtx, log, storage))
	router.Delete("/url", delete.New(myCtx, log, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout.Duration,
		WriteTimeout: cfg.HTTPServer.Timeout.Duration,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout.Duration,
	}

	log.Info("starting server", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}

func setupLogger(cfg *config.Logger) *slog.Logger {
	var log *slog.Logger

	var level slog.Level
	switch cfg.Level {
	case levelDebug:
		level = slog.LevelDebug
	case levelInfo:
		level = slog.LevelInfo
	}

	var out io.Writer
	switch cfg.Out {
	case outStd:
		out = os.Stdout
	}

	switch cfg.Type {
	case logText:
		log = slog.New(
			slog.NewTextHandler(out, &slog.HandlerOptions{
				Level: level,
			}),
		)
	case logJSON:
		log = slog.New(
			slog.NewJSONHandler(out, &slog.HandlerOptions{
				Level: level,
			}),
		)
	}

	return log
}
