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
	FindActiveSession(ctx context.Context, sid string) (bool, error)
}

func AuthCheck(next http.Handler, core Core, lg *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if errors.Is(err, http.ErrNoCookie) {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		result, err := core.FindActiveSession(r.Context(), session.Value)
		if err != nil {
			lg.Error("auth check error", "err", err.Error())
			next.ServeHTTP(w, r)
			return
		}

		lg.Info("HUI", result)

		if result == false {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, lg)
			return
		}

		//r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userId))

		next.ServeHTTP(w, r)
	})
}
