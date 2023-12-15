package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"restAPI/internal/config"
	"restAPI/internal/http-server/handlers/create"
	"restAPI/internal/http-server/handlers/deleteAll"
	"restAPI/internal/http-server/handlers/deleteById"
	"restAPI/internal/http-server/handlers/getAll"
	"restAPI/internal/http-server/handlers/getByDue"
	"restAPI/internal/http-server/handlers/getById"
	"restAPI/internal/http-server/handlers/getByTag"
	"restAPI/internal/storage/mapstorage"
	"restAPI/pkg/logger"
)

func main() {
	cfg := config.New()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting task store", slog.String("env", cfg.Env))
	log.Debug("debug is enabled")
	log.Debug("config value:", cfg)

	storage := mapstorage.NewTaskStore()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/task", create.New(log, storage))
	router.Get("/task", getAll.New(log, storage))
	router.Get("/task/{taskId}", getById.New(log, storage))
	router.Delete("/task", deleteAll.New(log, storage))
	router.Delete("/task/{taskId}", deleteById.New(log, storage))
	router.Get("/tag/{tag}", getByTag.New(log, storage))
	router.Get("/tag/{year}/{month}/{day}", getByDue.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start the server", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
	log.Error("Server stopped")
}
