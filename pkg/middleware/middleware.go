package middleware

import (
	"context"
	"errors"
	"filmoteka/pkg/models"
	httpResponse "filmoteka/pkg/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userId"

type Core interface {
	GetUserId(ctx context.Context, sid string) (uint64, error)
	GetRole(ctx context.Context, userId uint64) (string, error)
}

func AuthCheck(next http.Handler, core Core, lg *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if errors.Is(err, http.ErrNoCookie) {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		userId, err := core.GetUserId(r.Context(), session.Value)
		if err != nil {
			lg.Error("auth check error", "err", err.Error())
			next.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userId))
		if userId == 0 {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckRole(next http.Handler, core Core, lg *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, isAuth := r.Context().Value(UserIDKey).(uint64)
		if !isAuth {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		result, err := core.GetRole(r.Context(), userId)
		if err != nil {
			lg.Error("auth check error", "err", err.Error())
			next.ServeHTTP(w, r)
			return
		}

		if result != "admin" {
			response := models.Response{Status: http.StatusConflict, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		next.ServeHTTP(w, r)
	})
}
