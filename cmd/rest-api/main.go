package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/handlers/create"
	"restAPI/internal/http-server/handlers/deleteAll"
	"restAPI/internal/http-server/handlers/deleteById"
	"restAPI/internal/http-server/handlers/getAll"
	"restAPI/internal/http-server/handlers/getByDate"
	"restAPI/internal/http-server/handlers/getById"
	"restAPI/internal/http-server/handlers/getByTag"
	"restAPI/internal/repositories/tasks"

	"restAPI/internal/config"
	"restAPI/internal/db/postgres"
	"restAPI/pkg/logger"
)

func main() {
	cfg := config.New()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting task store", slog.String("env", cfg.Env))
	log.Debug("debug is enabled")
	log.Debug("config value:", cfg)

	storageCfg := postgres.NewConfig(log)
	log.Debug("storage config value:", storageCfg)
	conn, err := postgres.New(storageCfg)
	defer closeConn(conn)
	if err != nil {
		log.Error("cannot create ")
	}
	storage := tasks.NewRepository(conn, log)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/task", create.New(log, storage))
	router.Get("/task", getAll.New(log, storage))
	router.Get("/task/{taskId}", getById.New(log, storage))
	router.Delete("/task", deleteAll.New(log, storage))
	router.Delete("/task/{taskId}", deleteById.New(log, storage))
	router.Get("/tag/{tag}", getByTag.New(log, storage))
	router.Get("/due-date/{year}/{month}/{day}", getByDate.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start the server", slog.String("error", err.Error()))
	}
	log.Error("Server stopped")
}

func closeConn(conn *sqlx.DB) {
	err := conn.Close()
	if err != nil {
		panic("can't close connection to db")
	}
}
