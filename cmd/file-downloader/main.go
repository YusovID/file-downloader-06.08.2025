package main

import (
	"log/slog"
	"os"

	"github.com/YusovID/file-downloader-06.08.2025/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func init() {
	err := os.Setenv("CONFIG_PATH", "./config/local.yaml")
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting file-downloader", slog.String("env", cfg.Env))

	// TODO: init router

	// TODO: init server
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
