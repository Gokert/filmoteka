package middleware

import (
	"context"
	"errors"
	"filmoteka/pkg/models"
	httpResponse "filmoteka/pkg/response"
	core_profiles "filmoteka/usecase/profiles"
	core_session "filmoteka/usecase/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userId"

type Middleware struct {
	Lg       *logrus.Logger
	Sessions core_session.ISessions
	Profiles core_profiles.IProfiles
}

func (m *Middleware) AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if errors.Is(err, http.ErrNoCookie) {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, m.Lg)
			return
		}

		userId, err := m.Sessions.GetUserId(r.Context(), session.Value)
		if err != nil {
			m.Lg.Error("auth check error", "err", err.Error())
			next.ServeHTTP(w, r)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userId))
		if userId == 0 {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, m.Lg)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CheckRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, isAuth := r.Context().Value(UserIDKey).(uint64)
		if !isAuth {
			response := models.Response{Status: http.StatusUnauthorized, Body: nil}
			httpResponse.SendResponse(w, r, &response, m.Lg)
			return
		}

		result, err := m.Profiles.GetRole(r.Context(), userId)
		if err != nil {
			m.Lg.Error("auth check error", "err", err.Error())
			next.ServeHTTP(w, r)
			return
		}

		if result != "admin" {
			response := models.Response{Status: http.StatusConflict, Body: nil}
			httpResponse.SendResponse(w, r, &response, m.Lg)
			return
		}

		next.ServeHTTP(w, r)
	})
}
