package create

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"time"
)

type Request struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

type Response struct {
	response.Response
	TaskId int64 `json:"task_id"`
}

type TaskCreater interface {
	CreateTask(text string, tags []string, date time.Time) (int64, error)
}

func New(log *slog.Logger, taskCreater TaskCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		taskId, err := taskCreater.CreateTask(req.Text, req.Tags, time.Now())
		if err != nil {
			log.Error("failed to create task:", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, response.Error("failed to create task"))
			return
		}
		log.Info("task created", slog.Int64("taskID", taskId))
		render.JSON(w, r, Response{
			Response: response.Ok(),
			TaskId:   taskId,
		})

	}

}
