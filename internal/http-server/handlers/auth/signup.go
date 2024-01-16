package auth

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/internal/model"
	"restAPI/pkg/lib/verification"
)

type signUpRequest struct {
	model.User
}

type signUpResponse struct {
	UserID int64 `json:"user_id"`
}

//go:generate mockgen -source=signup.go -destination=mocks/signup-mock.go

type userCreater interface {
	CreateUser(user model.User) (int64, error)
}

// SignUp
// @Summary SignUp
// @Tags Authorization
// @Description Sign up handler
// @ID signUp
// @Accept json
// @Produce json
// @Param input body signUpRequest true "user info"
// @Success 201 {object} signUpResponse
// @Failure 400 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /auth/sign-up [post]
func SignUp(log *slog.Logger, creater userCreater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		//request decoding
		var req signUpRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, response.Message{
				Msg: "failed to decode request",
			})
			return
		}
		//user fields verification
		if !verification.User(req.User) {
			log.Error("wrong json fields")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "wrong json fields",
			})
			return
		}
		log.Debug("auth decoded", slog.String("firstname", req.FirstName), slog.String("secondname", req.SecondName), slog.String("login", req.Login))
		log.Info("request body decoded", slog.Any("request", req))

		//creating user and saving him in db
		//TODO:Различать разные ошибки, например: логин уже существует
		userId, err := creater.CreateUser(req.User)
		if err != nil {
			log.Error("failed to create auth:", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "failed to register",
			})
			return
		}
		log.Info("auth created", slog.Int64("userID", userId))

		//sending success response to client
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, signUpResponse{
			UserID: userId,
		})
	}
}
