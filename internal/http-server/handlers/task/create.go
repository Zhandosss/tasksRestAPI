package task

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

type CreateResponse struct {
	TaskId int64 `json:"task_id"`
}

type Creater interface {
	CreateTask(text string, tags []string, date time.Time, ownerID int64) (int64, error)
}

func Create(log *slog.Logger, creater Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		var req Request

		userID := r.Context().Value("userID").(int64)

		if userID == 0 {
			log.Error("couldn't get userID")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, "failed to get user id")
			return
		}

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		taskId, err := creater.CreateTask(req.Text, req.Tags, time.Now(), userID)
		if err != nil {
			log.Error("failed to create task:", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "failed to create task")
			return
		}
		log.Info("task created", slog.Int64("taskID", taskId))
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, CreateResponse{
			TaskId: taskId,
		})

	}

}
