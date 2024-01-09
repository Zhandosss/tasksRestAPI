package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"restAPI/internal/config"
	"restAPI/internal/db"
	"restAPI/internal/http-server/handlers/date"
	"restAPI/internal/http-server/handlers/tag"
	"restAPI/internal/http-server/handlers/task"
	"restAPI/internal/http-server/handlers/user"
	jwtAuth "restAPI/internal/http-server/middleware/JWTAuth"
	"restAPI/internal/repositories"
	"restAPI/internal/service"
	"restAPI/pkg/logger"
)

func main() {
	cfg := config.New()

	log := logger.SetupLogger(cfg.Env)

	log.Info("starting task store", slog.String("env", cfg.Env))
	log.Debug("debug is enabled")
	log.Debug("config value:", cfg)

	conn, err := db.New(&cfg.DB)
	if err != nil {
		log.Error("cannot create connection", slog.String("error", err.Error()))
	}
	defer closeConn(conn)

	rep := repositories.New(conn, log)

	services := service.New(rep)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/auth", func(router chi.Router) {
		router.Post("/sign-up", user.SignUp(log, services))
		router.Post("/sign-in", user.SignIn(log, services))
	})

	router.Group(func(router chi.Router) {
		router.Route("/tasks", func(router chi.Router) {
			router.Use(middleware.BasicAuth("restAPI", map[string]string{
				cfg.Admin.Login: cfg.Admin.Password,
			}))
			router.Delete("/", task.DeleteAll(log, services))
			router.Get("/", task.GetAll(log, services))

		})
	})

	router.Group(func(router chi.Router) {
		router.Use(jwtAuth.New(log, services))
		router.Route("/tasks", func(router chi.Router) {
			router.Post("/", task.Create(log, services))
			router.Get("/{taskId}", task.Get(log, services))
			router.Delete("/{taskId}", task.Delete(log, services))
		})
		router.Route("/tag", func(router chi.Router) {
			router.Get("/{tag}", tag.Get(log, services))
		})
		router.Route("/date", func(router chi.Router) {
			router.Get("/{year}/{month}/{day}", date.Get(log, services))

		})
	})

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
