package tag

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/model"
)

type Response struct {
	Tasks []model.Task
}

type GetterByTag interface {
	GetTasksByTag(tag string, userID int64) ([]model.Task, error)
}

func Get(log *slog.Logger, getterByTag GetterByTag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		userID := r.Context().Value("userID").(int64)

		if userID == 0 {
			log.Error("couldn't get userID")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, "failed to get user id")
			return
		}

		tag := chi.URLParam(r, "tag")
		if tag == "" {
			log.Error("there is no tag")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "there is no tag")
			return
		}
		tasks, err := getterByTag.GetTasksByTag(tag, userID)
		if err != nil {
			log.Error("get-by-tag", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, "couldn't find any data")
			return
		}
		log.Info("tasks copied by tag", slog.String("tag", tag))
		render.JSON(w, r, Response{
			Tasks: tasks,
		})
	}

}
