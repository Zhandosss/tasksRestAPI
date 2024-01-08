package signin

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Token string `json:"token"`
}

type tokenGenerator interface {
	GenerateToken(login, password string) (string, error)
}

func New(log *slog.Logger, generator tokenGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to decode request")
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		token, err := generator.GenerateToken(req.Login, req.Password)
		if err != nil {
			log.Error("generating token error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, "failed to login")
		}
		log.Info("token created")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, Response{
			Token: token,
		})
	}
}
