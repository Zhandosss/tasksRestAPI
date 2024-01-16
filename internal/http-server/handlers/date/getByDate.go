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

type getTaskResponse struct {
	Tasks []model.Task `json:"tasks"`
}

type getterByDate interface {
	GetTasksByDate(day, month, year int, userID int64) ([]model.Task, error)
}

// Get task by date
// @Summary Get
// @Security ApiKeyPath
// @Tags Date
// @Description Get user task by date
// @ID getTaskByDate
// @Param year path int true "year"
// @Param month path int true "month"
// @Param day path int true "day"
// @Produce json
// @Success 200 {object} getTaskResponse
// @Failure 400,401 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /date/{year}/{month}/{day} [get]
func Get(log *slog.Logger, getterByDate getterByDate) http.HandlerFunc {
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

		tasks, err := getterByDate.GetTasksByDate(dayInt, monthInt, yearInt, userID)
		if err != nil {
			log.Error("get tasks by date", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "couldn't find any tasks by date",
			})
			return
		}
		log.Info("tasks copied by data", slog.String("day", day), slog.String("month", month), slog.String("year", year))
		render.JSON(w, r, getTaskResponse{
			Tasks: tasks,
		})

	}
}
