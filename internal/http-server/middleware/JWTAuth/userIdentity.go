package jwtAuth

import (
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strings"
)

type tokenParser interface {
	ParseToken(inputToken string) (int64, error)
}

func New(log *slog.Logger, parser tokenParser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("JWTAuth middleware enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				log.Error("authorization header is empty")
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, "empty authorization header")
				return
			}
			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 {
				log.Error("authorization header is wrong")
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, "bad token")
				return
			}
			userID, err := parser.ParseToken(headerParts[1])
			if err != nil {
				log.Error("couldn't parse token", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, "bad token")
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			newReq := r.WithContext(ctx)
			next.ServeHTTP(w, newReq)
		}
		return http.HandlerFunc(fn)
	}
}
