package deleteAll

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
)

type AllDeleter interface {
	DeleteAllTasks() error
}

func New(log *slog.Logger, allDeleter AllDeleter) http.HandlerFunc {
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
			render.JSON(w, r, response.Error("failed to delete all tasks"))
			return
		}
		log.Info("all tasks deleted")
		render.JSON(w, r, response.Ok())
	}
}
