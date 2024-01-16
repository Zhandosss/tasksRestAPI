package tag

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/model"
)

type getTaskResponse struct {
	Tasks []model.Task `json:"tasks"`
}

type getterByTag interface {
	GetTasksByTag(tag string, userID int64) ([]model.Task, error)
}

// Get task by tag
// @Summary Get
// @Security ApiKeyPath
// @Tags Tag
// @Description Get user task by tag
// @ID getTaskByTag
// @Param tag path string true "tag"
// @Produce json
// @Success 200 {object} getTaskResponse
// @Failure 400,403 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /tag/{tag} [get]
func Get(log *slog.Logger, getterByTag getterByTag) http.HandlerFunc {
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

		tag := chi.URLParam(r, "tag")
		if tag == "" {
			log.Error("there is no tag")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "failed to get tag from url",
			})
			return
		}
		tasks, err := getterByTag.GetTasksByTag(tag, userID)
		if err != nil {
			log.Error("couldn't find any task by tag", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "couldn't find any task by tag",
			})
			return
		}
		log.Info("tasks copied by tag", slog.String("tag", tag))
		render.JSON(w, r, getTaskResponse{
			Tasks: tasks,
		})
	}

}
