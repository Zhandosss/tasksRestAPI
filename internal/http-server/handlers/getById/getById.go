package getById

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/model"
	"strconv"
)

type Response struct {
	response.Response
	Task model.Task `json:"task"`
}

type TaskGetter interface {
	GetTask(taskID int64) (model.Task, error)
}

func New(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		taskIdString := chi.URLParam(r, "taskId")
		if taskIdString == "" {
			log.Error("failed to get task id from url")
			render.JSON(w, r, response.Error("failed to get task id from url"))
			return
		}
		taskId, err := strconv.Atoi(taskIdString)
		if err != nil {
			log.Error("incorrect task id record", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, response.Error("incorrect task id record"))
			return
		}
		task, err := taskGetter.GetTask(int64(taskId))
		if err != nil {
			log.Error("can't found task", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, response.Error("can't found task"))
			return
		}
		log.Info("task copied by id", slog.Int("id", taskId))
		render.JSON(w, r, Response{
			Response: response.Ok(),
			Task:     task,
		})
	}
}
