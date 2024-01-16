package task

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

type getTaskResponse struct {
	Task model.Task `json:"task"`
}

type taskGetter interface {
	GetTask(taskID, UserID int64) (model.Task, error)
}

// Get task by ID
// @Summary Get
// @Security ApiKeyPath
// @Tags Task
// @Description Get user task by ID
// @ID getTaskByID
// @Param task_id path int true "task ID"
// @Produce json
// @Success 200 {object} getTaskResponse
// @Failure 400,403 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /tasks/{taskId} [get]
func Get(log *slog.Logger, taskGetter taskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		userID := r.Context().Value("userID").(int64)

		if userID <= 0 {
			log.Error("couldn't get userID")
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, response.Message{
				Msg: "failed to get auth id",
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
		taskId, err := strconv.Atoi(taskIdString)
		if err != nil {
			log.Error("incorrect task id record", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "incorrect task id record",
			})
			return
		}
		task, err := taskGetter.GetTask(int64(taskId), userID)
		if err != nil {
			log.Error("can't found task", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "can't found task",
			})
			return
		}
		log.Info("task copied by id", slog.Int("id", taskId))
		render.JSON(w, r, getTaskResponse{
			Task: task,
		})
	}
}
