package user

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/model"
)

type SignUpRequest struct {
	model.User
}

type SignUpResponse struct {
	UserID int64 `json:"user_id"`
}

type Creater interface {
	CreateUser(user model.User) (int64, error)
}

func SignUp(log *slog.Logger, creater Creater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		var req SignUpRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}
		log.Debug("user decoded", slog.String("firstname", req.FirstName), slog.String("secondname", req.SecondName), slog.String("login", req.Login))
		log.Info("request body decoded", slog.Any("request", req))

		userId, err := creater.CreateUser(req.User)
		if err != nil {
			log.Error("failed to create user:", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to register")
			return
		}
		log.Info("user created", slog.Int64("userID", userId))
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, SignUpResponse{
			UserID: userId,
		})
	}
}
