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

type taskDeleter interface {
	DeleteTask(taskID, userID int64) error
}

// Delete task by ID
// @Summary Delete
// @Security ApiKeyPath
// @Tags Task
// @Description Delete user task by ID
// @ID deleteTaskByID
// @Produce json
// @Param task_id path int true "task ID"
// @Success 204
// @Failure 400,401,404 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /tasks/{taskId} [delete]
func Delete(log *slog.Logger, deleter taskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		userID := r.Context().Value("userID").(int64)

		if userID <= 0 {
			log.Error("couldn't get userID")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, response.Message{
				Msg: "failed to get auth id",
			})
			return
		}

		taskIDString := chi.URLParam(r, "taskId")
		if taskIDString == "" {
			log.Error("failed to get task id from url")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "failed to get task id from url",
			})
			return
		}
		taskID, err := strconv.Atoi(taskIDString)
		if err != nil {
			log.Error("incorrect task id record", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "incorrect task id record",
			})
			return
		}
		err = deleter.DeleteTask(int64(taskID), userID)
		if errors.Is(err, repositories.ErrNoTask) {
			log.Error("there is no task", slog.Int64("userID", userID), slog.Int("taskID", taskID))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Message{
				Msg: "there no task with this taskID",
			})
			return
		}
		if err != nil {
			log.Error("can't found task", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "can't found task",
			})
			return
		}
		log.Info("task was deleted by Id", slog.Int("id", taskID))
		w.WriteHeader(http.StatusNoContent)
	}
}
