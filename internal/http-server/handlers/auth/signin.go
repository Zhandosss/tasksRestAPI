package auth

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"restAPI/internal/http-server/response"
	"restAPI/pkg/lib/verification"
)

type signInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

type tokenGenerator interface {
	GenerateToken(login, password string) (string, error)
}

// SignIn
// @Summary SignIn
// @Tags Authorization
// @Description Login handler
// @ID login
// @Accept json
// @Produce json
// @Param input body signInRequest true "Login and password"
// @Success 201 {object} signInResponse
// @Failure 400 {object} response.Message
// @Failure 500 {object} response.Message
// @Failure default {object} response.Message
// @Router /auth/sign-in [post]
func SignIn(log *slog.Logger, generator tokenGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		//request decoding
		var req signInRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "failed to decode request",
			})
			return
		}

		//login and password verification
		if !verification.LoginAndPassword(req.Login, req.Password) {
			log.Error("wrong json fields")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Message{
				Msg: "wrong json fields",
			})
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		//token generation
		token, err := generator.GenerateToken(req.Login, req.Password)
		if err != nil {
			log.Error("generating token error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Message{
				Msg: "failed to login",
			})
			return
		}
		//sending success response to client
		log.Info("token created")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, signInResponse{
			Token: token,
		})
	}
}
