package task

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/model"
	"restAPI/pkg/lib/verification"
	"time"
)

type CreateRequest struct {
	model.Task
}

type CreateResponse struct {
	TaskID int64 `json:"task_id"`
}

type Creater interface {
	CreateTask(task model.Task) (int64, error)
}

func Create(log *slog.Logger, creater Creater, callTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		var req CreateRequest
		//getting user id from request context
		userID := r.Context().Value("userID").(int64)

		if userID <= 0 {
			log.Error("incorrect userID", slog.Int64("userID", userID))
			w.WriteHeader(http.StatusForbidden)
			render.JSON(w, r, response.Message{
				Msg: "incorrect userID",
			})
			return
		}

		//decoding request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "failed to decode request",
			})
			return
		}

		req.Task.OwnerID = userID
		//task verification
		//TODO: ADD request body in logger output
		if !verification.Task(req.Task, callTime) {
			log.Error("incorrect task information")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "incorrect task information",
			})
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		taskId, err := creater.CreateTask(req.Task)
		if err != nil {
			log.Error("failed to create task:", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "failed to create task",
			})
			return
		}
		log.Info("task created", slog.Int64("taskID", taskId))

		//success request
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, CreateResponse{
			TaskID: taskId,
		})

	}

}
