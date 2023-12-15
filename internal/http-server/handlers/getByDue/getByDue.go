package getByDue

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/model"
	"time"
)

type Response struct {
	response.Response
	Tasks []model.Task
}

type GetterByDue interface {
	GetTasksByDueDate(year int, month time.Month, day int) ([]model.Task, error)
}

func New(log *slog.Logger, getterByDue GetterByDue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		year := chi.URLParam(r, "year")
		month := chi.URLParam(r, "month")
		day := chi.URLParam(r, "day")
	}
}
