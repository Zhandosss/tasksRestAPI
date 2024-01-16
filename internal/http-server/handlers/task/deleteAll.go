package task

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
)

type allTasksDeleter interface {
	DeleteAllByUser(userID int64) error
}

// DeleteAll user tasks
// @Summary DeleteAll
// @Security ApiKeyPath
// @Tags Task
// @Description Delete all user tasks
// @ID deleteAllUserTasks
// @Produce json
// @Success 204
// @Failure 401 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /tasks/ [delete]
func DeleteAll(log *slog.Logger, deleter allTasksDeleter) http.HandlerFunc {
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

		err := deleter.DeleteAllByUser(userID)

		if err != nil {
			log.Error("couldn't delete all tasks by this user", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "couldn't delete all task",
			})
			return
		}

		log.Info("all user tasks deleted", slog.Int64("userID", userID))
		w.WriteHeader(http.StatusNoContent)
	}
}
