package main

import (
	"log/slog"
	"os"
	"url_shortener/internal/config"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//TODO: init config: cleanenv +
	cfg := config.MustLoad()

	//TODO: init logger: slog +
	log := setupLogger(cfg.Env)
	log.Info("Starting server", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	//TODO: init storege: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	//TODO: init router: chi, "chi render"

	//TODO: start server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
