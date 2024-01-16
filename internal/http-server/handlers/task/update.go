package task

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/repositories"
	"strconv"
)

type updateRequest struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

type taskUpdater interface {
	UpdateTask(taskID, userID int64, Text string, Tags []string) error
}

// Update task by ID
// @Summary Update
// @Security ApiKeyPath
// @Tags Task
// @Description Update user task by ID
// @ID updateTaskByID
// @Param input body updateRequest true "new text and tags"
// @Param task_id path int true "task ID"
// @Produce json
// @Success 204
// @Failure 400,401,404 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /tasks/{taskId} [put]
func Update(log *slog.Logger, updater taskUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		var req updateRequest

		userID := r.Context().Value("userID").(int64)

		if userID <= 0 {
			log.Error("couldn't get userID")
			w.WriteHeader(http.StatusUnauthorized)
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
		if errors.Is(err, repositories.ErrNoTask) {
			log.Error("there is no task", slog.Int64("userID", userID), slog.Int("taskID", taskID))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Message{
				Msg: "there no task with this taskID",
			})
			return
		}
		if err != nil {
			log.Error("can't update task", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "can't update task",
			})
			return
		}
		log.Info("task updated")
		w.WriteHeader(http.StatusNoContent)
	}
}
