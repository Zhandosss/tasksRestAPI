package get_all

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/model"
	"restAPI/internal/repositories"
)

type Response struct {
	AllTasks []model.Task `json:"all_tasks"`
}

type AllGetter interface {
	GetAllTasks() ([]model.Task, error)
}

func New(log *slog.Logger, allGetter AllGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		ans, err := allGetter.GetAllTasks()
		if err != nil {
			log.Error("failed to get all tasks", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			if errors.Is(err, repositories.ErrEmptyTable) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			render.JSON(w, r, "failed to get all tasks")
			return
		}
		log.Info("all tasks was copied")
		render.JSON(w, r, Response{
			AllTasks: ans,
		})
	}

}
