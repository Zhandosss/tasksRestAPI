package task

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"strconv"
)

type UpdateRequest struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

type Updater interface {
	UpdateTask(taskID, userID int64, Text string, Tags []string) error
}

func Update(log *slog.Logger, updater Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		var req UpdateRequest

		userID := r.Context().Value("userID").(int64)

		if userID <= 0 {
			log.Error("couldn't get userID")
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, response.Message{
				Msg: "failed to get auth id",
			})
			return
		}

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "failed to decode request",
			})
			return
		}

		if req.Text == "" {
			log.Error("there is no text in updated version")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "there is no text in task",
			})
			return
		}

		taskIdString := chi.URLParam(r, "taskId")
		if taskIdString == "" {
			log.Error("failed to get task id from url")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "failed to get task id from url",
			})
			return
		}
		taskID, err := strconv.Atoi(taskIdString)
		if err != nil {
			log.Error("incorrect task id record", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "incorrect task id record",
			})
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		err = updater.UpdateTask(int64(taskID), userID, req.Text, req.Tags)

		if err != nil {
			log.Error("can't update task", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Message{
				Msg: "can't update task",
			})
			return
		}
		log.Info("task updated")
		render.JSON(w, r, response.Message{
			Msg: "task updated",
		})
	}
}
