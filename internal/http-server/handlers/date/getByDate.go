package date

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
	Tasks []model.Task `json:"tasks"`
}

type GetterByDue interface {
	GetTasksByDate(day, month, year int, userID int64) ([]model.Task, error)
}

func Get(log *slog.Logger, getterByDue GetterByDue) http.HandlerFunc {
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

		year := chi.URLParam(r, "year")
		month := chi.URLParam(r, "month")
		day := chi.URLParam(r, "day")
		if !verification.Date(day, month, year) {
			log.Error(`incorrect data format`, slog.String("data", fmt.Sprintf("%s:%s:%s", day, month, year)))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "incorrect data format",
			})
			return
		}
		dayInt, _ := strconv.Atoi(day)
		monthInt, _ := strconv.Atoi(month)
		yearInt, _ := strconv.Atoi(year)

		tasks, err := getterByDue.GetTasksByDate(dayInt, monthInt, yearInt, userID)
		if err != nil {
			log.Error("get tasks by date", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Message{
				Msg: "couldn't find any tasks by date",
			})
			return
		}
		log.Info("tasks copied by data", slog.String("day", day), slog.String("month", month), slog.String("year", year))
		render.JSON(w, r, Response{
			Tasks: tasks,
		})

	}
}
