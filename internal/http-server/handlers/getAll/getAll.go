package getAll

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/model"
)

type Response struct {
	response.Response
	AllTasks []model.Task `json:"all_tasks"`
}

type AllGetter interface {
	GetAllTasks() ([]model.Task, error)
}

func New(log *slog.Logger, allGetter AllGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		ans, err := allGetter.GetAllTasks()
		if err != nil {
			log.Error("failed to get all tasks", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, response.Error("failed to get all tasks"))
			return
		}
		log.Info("all tasks was copied")
		render.JSON(w, r, Response{
			Response: response.Ok(),
			AllTasks: ans,
		})
	}

}
