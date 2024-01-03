package deleteById

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type TaskDeleter interface {
	DeleteTask(taskID int64) error
}

func New(log *slog.Logger, deleter TaskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		taskIdString := chi.URLParam(r, "taskId")
		if taskIdString == "" {
			log.Error("failed to get task id from url")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to get task id from url")
			return
		}
		taskId, err := strconv.Atoi(taskIdString)
		if err != nil {
			log.Error("incorrect task id record", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "incorrect task id record")
			return
		}
		err = deleter.DeleteTask(int64(taskId))
		if err != nil {
			log.Error("can't found task", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, "can't found task")
			return
		}
		log.Info("task was deleted by Id", slog.Int("id", taskId))
		w.WriteHeader(http.StatusNoContent)
	}
}
