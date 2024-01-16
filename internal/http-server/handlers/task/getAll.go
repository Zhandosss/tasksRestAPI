package task

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/model"
)

type getAllResponse struct {
	Tasks []model.Task `json:"tasks"`
}

type allGetterByUser interface {
	GetAllByUser(userID int64) ([]model.Task, error)
}

// GetAll user tasks
// @Summary GetAll
// @Security ApiKeyPath
// @Tags Task
// @Description Get all user tasks
// @ID getAllUserTasks
// @Produce json
// @Success 200 {object} getAllResponse
// @Failure 403 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /tasks/ [get]
func GetAll(log *slog.Logger, getter allGetterByUser) http.HandlerFunc {
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

		tasks, err := getter.GetAllByUser(userID)
		if err != nil {
			log.Error("couldn't get all tasks by this user", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "couldn't get all task",
			})
			return
		}

		log.Info("all user tasks copied", slog.Int64("userID", userID))
		render.JSON(w, r, getAllResponse{
			Tasks: tasks,
		})
	}
}
