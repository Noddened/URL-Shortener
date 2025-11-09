package main

import (
	"log/slog"
	"os"

	"github.com/Noddened/URL-Shortener/internal/config"
	"github.com/Noddened/URL-Shortener/internal/lib/logger/sl"
	"github.com/Noddened/URL-Shortener/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) //Каждое сообщение будет также содержать инфу об окружении

	log.Info("initializing server", slog.String("address", cfg.Address)) //Также выводим сообщение с адресом
	log.Debug("logger debug mode enabled")                               // тут итак понятно

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failde to init storage", sl.Err(err))
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID) // request_id в каждом запросе для трейсинга
	router.Use(middleware.Logger)    // логируем все запросы
	router.Use(middleware.Recoverer) // если внутри обработчика запросов произойдет паника, то сервер не упадет
	// упал -> отжался -> встал
	router.Use(mwLogger.New(log))
	router.Use(middleware.URLFormat) // парсер url поступающих запросов
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
