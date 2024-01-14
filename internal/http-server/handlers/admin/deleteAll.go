package admin

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/repositories"
)

type AllDeleter interface {
	DeleteAllTasks() error
}

func DeleteAll(log *slog.Logger, allDeleter AllDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("request", middleware.GetReqID(r.Context())),
		)

		err := allDeleter.DeleteAllTasks()
		if err != nil {
			log.Error("failed to delete all tasks", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			if errors.Is(err, repositories.ErrEmptyTable) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			render.JSON(w, r, "failed to delete all tasks")
			return
		}
		log.Info("all tasks deleted")
		w.WriteHeader(http.StatusNoContent)
	}
}
