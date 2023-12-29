package getByDate

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/model"
	"restAPI/pkg/lib/verification"
	"strconv"
)

type Response struct {
	response.Response
	Tasks []model.Task
}

type GetterByDue interface {
	GetTasksByDate(day, month, year int) ([]model.Task, error)
}

func New(log *slog.Logger, getterByDue GetterByDue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		year := chi.URLParam(r, "year")
		month := chi.URLParam(r, "month")
		day := chi.URLParam(r, "day")
		if !verification.Date(day, month, year) {
			log.Error(`incorrect data format`, slog.String("data", fmt.Sprintf("%s:%s:%s", day, month, year)))
			render.JSON(w, r, response.Error("incorrect data format"))
			return
		}
		dayInt, _ := strconv.Atoi(day)
		monthInt, _ := strconv.Atoi(month)
		yearInt, _ := strconv.Atoi(year)

		tasks, err := getterByDue.GetTasksByDate(dayInt, monthInt, yearInt)
		if err != nil {
			log.Error("get tasks by date", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, response.Error("couldn't find any data"))
			return
		}
		log.Info("tasks copied by data", slog.String("day", day), slog.String("month", month), slog.String("year", year))
		render.JSON(w, r, Response{
			Response: response.Ok(),
			Tasks:    tasks,
		})

	}
}
