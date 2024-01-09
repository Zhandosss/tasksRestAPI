package getById

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/model"
	"strconv"
)

type Response struct {
	Task model.Task `json:"task"`
}

type TaskGetter interface {
	GetTask(taskID, UserID int64) (model.Task, error)
}

func New(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
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
		task, err := taskGetter.GetTask(int64(taskId), userID)
		if err != nil {
			log.Error("can't found task", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, "can't found task")
			return
		}
		log.Info("task copied by id", slog.Int("id", taskId))
		render.JSON(w, r, Response{
			Task: task,
		})
	}
}
